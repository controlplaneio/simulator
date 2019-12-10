#!/bin/bash

git clone https://github.com/controlplaneio/kubesec-webhook.git webhook

cd webhook

apt update && apt install -y make

make certs

make deploy

cd ..

rm -rf webhook

kubectl create namespace master-encirclement

kubectl label namespaces master-encirclement kubesec-validation=enabled

sed -i '18i\ \ \ \ \- --disable-admission-plugins=MutatingAdmissionWebhook' /etc/kubernetes/manifests/kube-apiserver.yaml

apt remove -y make 
