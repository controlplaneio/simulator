#!/bin/bash
set -Eeuxo pipefail

sudo /tmp/authorized_keys.sh sublimino 06kellyjac rowan-baker wakeward jpts
rm /tmp/authorized_keys.sh

# Install necessary dependencies
sudo apt-get update -y
sudo DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" dist-upgrade
sudo apt-get update
sudo apt-get -y -qq install curl wget git vim apt-transport-https ca-certificates cloud-init figlet

# disable auto-updates
sudo apt purge --auto-remove -y -qq unattended-upgrades
sudo systemctl disable apt-daily-upgrade.timer
sudo systemctl mask apt-daily-upgrade.service
sudo systemctl disable apt-daily.timer
sudo systemctl mask apt-daily.service

# Temporarily hold these packages back due to a bug:
# https://github.com/kubernetes/release/issues/2866
# If crictl is forced up to 1.26.x, it expects the new CRI runtiem api and kubeadm fails
VERSION='1.24.9*'
CRICTL_VERSION='1.25.*'
cat <<EOF | sudo bash
set -Eeuxo pipefail

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main"  > /etc/apt/sources.list.d/kubernetes.list

apt update
apt install -y --allow-downgrades \
  kubelet=${VERSION} kubeadm=${VERSION} kubectl=${VERSION} \
  containerd \
  awscli \
  jq \
  cri-tools=${CRICTL_VERSION}

kubeadm config images pull &

mkdir -p /run/download
wget https://github.com/mikefarah/yq/releases/download/3.4.1/yq_linux_amd64 -O /run/download/yq
install /run/download/yq /usr/bin

mkdir -p /etc/containerd
containerd config default > /etc/containerd/config.toml
sed -i "s/SystemdCgroup = false/SystemdCgroup = true/g" /etc/containerd/config.toml
sed -i -e "s/disable_apparmor.*$/disable_apparmor = true/g" /etc/containerd/config.toml
systemctl restart containerd

echo "runtime-endpoint: unix:///run/containerd/containerd.sock" > /etc/crictl.yaml

systemctl enable containerd
systemctl stop kubelet

apt clean

wait

EOF
