#!/usr/bin/env bash

mkdir -p /etc/containerd/certs.d/reg.rescue.drop
sed -i -e  's|config_path = ""|config_path = "/etc/containerd/certs.d"|' /etc/containerd/config.toml

cat <<EOF > /etc/containerd/certs.d/reg.rescue.drop/hosts.toml
server = "http://__REGISTRY_IP__:30080"

[host."http://__REGISTRY_IP__:30080"]
  capabilities = ["pull", "resolve", "push"]
  skip_verify = true
EOF

systemctl restart containerd
