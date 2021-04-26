#!/bin/bash
set -Eeuxo pipefail

sudo /tmp/authorized_keys.sh sublimino denhamparry 06kellyjac ipedrazas rowan-baker wakeward
rm /tmp/authorized_keys.sh

# Install necessary dependencies
sudo apt-get update -y
sudo DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" dist-upgrade
sudo apt-get update
sudo apt-get -y -qq install curl wget git vim apt-transport-https ca-certificates

VERSION='1.20.*'
cat <<EOF | sudo bash
add-apt-repository --yes ppa:rmescandon/yq
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main"  > /etc/apt/sources.list.d/kubernetes.list

mkdir /run/download
apt update
apt install -y \
  kubelet=${VERSION} kubeadm=${VERSION} kubectl=${VERSION} \
  docker.io \
  awscli \
  jq \
  yq=3*

kubeadm config images pull &

wget https://github.com/kubernetes-incubator/cri-tools/releases/download/v1.11.1/crictl-v1.11.1-linux-amd64.tar.gz -O /run/download/crictl.tgz
tar -C /usr/bin -xzf /run/download/crictl.tgz
chmod 754 /usr/bin/crictl
chown root:root /usr/bin/crictl

systemctl enable docker
systemctl daemon-reload
systemctl restart docker
systemctl stop kubelet

wait

EOF
