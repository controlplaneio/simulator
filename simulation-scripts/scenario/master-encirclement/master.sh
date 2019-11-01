#!/bin/bash

git clone https://github.com/controlplaneio/kubesec-webhook.git webhook &> /dev/null

cd webhook

apt update && apt install -y make &> /dev/null

make certs &> /dev/null

make deploy &> /dev/null

cd ..

rm -rf webhook  &> /dev/null

kubectl create namespace master-encirclement &> /dev/null

kubectl label namespaces master-encirclement kubesec-validation=enabled &> /dev/null

sed -i '18i\ \ \ \ \- --disable-admission-plugins=MutatingAdmissionWebhook' /etc/kubernetes/manifests/kube-apiserver.yaml

apt remove -y make  &> /dev/null
