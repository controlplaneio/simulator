#!/bin/bash

set -Eeuo pipefail

export KUBECONFIG=/etc/kubernetes/admin.conf

openssl genrsa -out player.key 2048
openssl req -new -key player.key -out player.csr -subj "/CN=player"

# We deliberately do a replace here to force deletion of the old resource due to wierd behavior in the CSR api. See JCP for more info.
kubectl replace --force -f - <<-EOF
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: player
spec:
  request: $(base64 -w0 player.csr)
  usages: ['digital signature', 'key encipherment', 'client auth']
  signerName: kubernetes.io/kube-apiserver-client
EOF

## double check the right CSR was used
mod1=$(kubectl get csr player -o jsonpath='{.spec.request}' | base64 --decode | openssl req -noout -modulus -in -)
mod2=$(openssl req -noout -modulus -in player.csr)
if [[ "$mod1" != "$mod2" ]]; then
    echo "CSR was not set correctly, aborting."
    exit 1
fi

kubectl certificate approve player
kubectl wait --for=condition=Approved csr player
# Wait to be issued
sleep 3

kubectl get csr player -o jsonpath='{.status.certificate}' | base64 --decode > player.crt

export KUBECONFIG=./player.kubeconfig

SERVER=$(awk '/server:/ {print $2}' /etc/kubernetes/kubelet.conf)
kubectl config set-cluster default --certificate-authority=/etc/kubernetes/pki/ca.crt --server="$SERVER"
kubectl config set-credentials player --client-key player.key --client-certificate player.crt --embed-certs
kubectl config set-context player --cluster default --user player
kubectl config use-context player

export KUBECONFIG=/etc/kubernetes/admin.conf
kubectl -n kube-public create cm player-config --from-file=kubeconfig=<(base64 -w0 player.kubeconfig) --dry-run=client -oyaml | kubectl apply --force -f -

kubectl apply -f - <<-EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  creationTimestamp: null
  name: read-public-cm
  namespace: kube-public
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
EOF
kubectl apply -f - <<-EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  creationTimestamp: null
  name: node-read-cm
  namespace: kube-public
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: read-public-cm
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: system:nodes
EOF
