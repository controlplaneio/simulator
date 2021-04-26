#!/bin/bash

set -Eeuxo pipefail

add-apt-repository --yes ppa:rmescandon/yq
apt update
apt install -y jq 'yq=3*'

yq r -j /etc/kubernetes/manifests/etcd.yaml | jq --arg auth "--client-cert-auth=true" '(.spec.containers[].command[] | select(. == $auth)) = "--client-cert-auth=false"' | yq r - > etcd.yaml
sed -i -E -e 's#, "--peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt"##g' -e 's#, "--trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt"##g' etcd.yaml
mv etcd.yaml /etc/kubernetes/manifests/etcd.yaml
