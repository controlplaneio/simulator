#!/bin/bash

set -Eeuxo pipefail

add-apt-repository --yes ppa:rmescandon/yq
apt update
apt install -y jq 'yq=3*'

yq w -i /var/lib/kubelet/config.yaml authentication.anonymous.enabled true
yq w -i /var/lib/kubelet/config.yaml authorization.mode AlwaysAllow

systemctl daemon-reload
systemctl restart kubelet.service
