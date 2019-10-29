#!/bin/bash

add-apt-repository ppa:rmescandon/yq
apt update
apt install -y yq
apt-get install -y jq

yq r -j /etc/kubernetes/manifests/kube-apiserver.yaml | jq --arg auth "--insecure-port=0" '(.spec.containers[].command[] | select(. == $auth)) = "--insecure-port=8080"' | jq '.spec.containers[].command += ["--insecure-bind-address=0.0.0.0"]' | yq r - > kube-apiserver.yaml
mv kube-apiserver.yaml /etc/kubernetes/manifests/kube-apiserver.yaml
