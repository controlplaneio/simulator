#!/bin/bash

add-apt-repository ppa:rmescandon/yq
apt update
apt install yq -y
apt-get install -y jq

#this doesn't work
yq r -j /etc/kubernetes/manifests/etcd.yaml | jq -r '.spec.containers[].command[]|select(. == "--client-cert-auth=true")'|="--client-cert-auth=false"

#write to replace "--client-cert-auth=true" with "--client-cert-auth=false"
