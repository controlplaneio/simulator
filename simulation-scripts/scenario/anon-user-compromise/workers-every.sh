#!/bin/bash

set -Euxo pipefail

add-apt-repository --yes ppa:rmescandon/yq
apt update
apt install 'yq=3*' -y

yq w -i /var/lib/kubelet/config.yaml authentication.anonymous.enabled true
yq w -i /var/lib/kubelet/config.yaml authorization.mode AlwaysAllow

systemctl daemon-reload
systemctl restart kubelet.service
