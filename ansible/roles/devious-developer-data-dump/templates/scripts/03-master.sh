#!/bin/bash

set -Eeuo pipefail

export KUBECONFIG=/etc/kubernetes/admin.conf

openssl genrsa -out runner.key 2048
openssl req -new -key runner.key -out runner.csr -subj "/CN=runner"

# We deliberately do a replace here to force deletion of the old resource due to wierd behavior in the CSR api. See JCP for more info.
kubectl replace --force -f - <<-EOF
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: runner
spec:
  request: $(base64 -w0 runner.csr)
  usages: ['digital signature', 'key encipherment', 'client auth']
  signerName: kubernetes.io/kube-apiserver-client
EOF

## double check the right CSR was used
mod1=$(kubectl get csr runner -o jsonpath='{.spec.request}' | base64 --decode | openssl req -noout -modulus -in -)
mod2=$(openssl req -noout -modulus -in runner.csr)
if [[ "$mod1" != "$mod2" ]]; then
    echo "CSR was not set correctly, aborting."
    exit 1
fi

kubectl certificate approve runner
kubectl wait --for=condition=Approved csr runner
# Wait to be issued
sleep 3

kubectl get csr runner -o jsonpath='{.status.certificate}' | base64 --decode > runner.crt

export KUBECONFIG=/etc/kubernetes/runner.conf
SERVER=$(awk '/server:/ {print $2}' /etc/kubernetes/kubelet.conf)
kubectl config set-cluster default --certificate-authority=/etc/kubernetes/pki/ca.crt --server="$SERVER" --embed-certs
kubectl config set-credentials runner --client-key runner.key --client-certificate runner.crt --embed-certs
kubectl config set-context runner --cluster default --user runner
kubectl config use-context runner
