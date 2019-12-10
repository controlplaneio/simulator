#! /bin/bash

add-apt-repository ppa:rmescandon/yq
apt update
apt install yq -y

yq w -i /var/lib/kubelet/config.yaml authentication.anonymous.enabled true
yq w -i /var/lib/kubelet/config.yaml authorization.mode AlwaysAllow

systemctl daemon-reload
systemctl restart kubelet.service
