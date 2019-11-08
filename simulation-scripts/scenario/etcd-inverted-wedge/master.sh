#!/bin/bash

add-apt-repository ppa:rmescandon/yq &> /dev/null
apt update &> /dev/null
apt install -y jq yq &> /dev/null

yq r -j /etc/kubernetes/manifests/etcd.yaml | jq --arg auth "--client-cert-auth=true" '(.spec.containers[].command[] | select(. == $auth)) = "--client-cert-auth=false"' | yq r - > etcd.yaml
mv etcd.yaml /etc/kubernetes/manifests/etcd.yaml
