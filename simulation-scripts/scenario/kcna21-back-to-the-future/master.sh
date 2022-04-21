#!/bin/bash

openssl genrsa -out player.key 2048
openssl req -new -key player.key -out player.csr -subj "/CN=player"
kubectl apply -f - <<EOF
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: player
spec:
  request: $(base64 -w0 player.csr)
  usages: ['digital signature', 'key encipherment', 'client auth']
EOF
kubectl certificate approve player
kubectl wait --for=condition=Approved csr player
sleep 5
kubectl get csr player -o jsonpath='{.status.certificate}' | base64 --decode > player.crt
kubectl config set-credentials player --client-key player.key --client-certificate player.crt --embed-certs
kubectl config set-context player --cluster kubernetes --user player
kubectl config use-context player
kubectl config delete-context kubernetes-admin@kubernetes
kubectl config delete-user kubernetes-admin
cp /root/.kube/config /home/ubuntu/.kube/config
mkdir -p ~ubuntu/.kube
chown ubuntu:ubuntu -R ~ubuntu/.kube/

timedatectl set-ntp off
date --set="25 OCT 1985 07:53:00"
echo "ubuntu ALL=(ALL) /usr/sbin/hwclock *" > /etc/sudoers.d/time
gpasswd -d ubuntu sudo
rm /etc/sudoers.d/90-cloud-init-users
chmod 755 /etc/sudoers.d
