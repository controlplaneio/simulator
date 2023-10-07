#!/usr/bin/env bash

set -Eeuo pipefail

curl -sSL 'https://get.helm.sh/helm-v3.12.1-linux-amd64.tar.gz' | tar xzf - --strip-components=1 linux-amd64/helm \
 && install -Dm755 helm /usr/local/bin/helm

IPNET=$(ip a s eth0 | awk '/inet / {print $2}')
IP=${IPNET%"/24"}

# Encapsulated Routing Mode /w performance tweaks
cat <<EOF > values.yaml
hubble:
  enabled: false
kubeProxyReplacement: strict
bpf:
  masquerade: true
  tproxy: true
bandwidthManager:
  enabled: true
endpointStatus:
  enabled: true
  status: policy
operator:
  replicas: 1
k8sServiceHost: $IP
k8sServicePort: 6443
EOF

helm upgrade --install --repo https://helm.cilium.io cilium cilium --version 1.13.4 --values ./values.yaml --namespace kube-system

kubectl wait -n kube-system --selector='app.kubernetes.io/name=cilium-agent' --for=condition=Ready --timeout=10m pods
