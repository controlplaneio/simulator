#!/bin/bash
set -Eeuxo pipefail

sudo /tmp/authorized_keys.sh sublimino denhamparry 06kellyjac rowan-baker wakeward jpts
rm /tmp/authorized_keys.sh

# Install necessary dependencies
sudo apt-get update -y
sudo DEBIAN_FRONTEND=noninteractive apt-get -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" dist-upgrade
sudo apt-get update
sudo apt-get -y -qq install curl wget git vim apt-transport-https ca-certificates cloud-init

VERSION='1.20.*'
cat <<EOF | sudo bash
set -Eeuxo pipefail

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main"  > /etc/apt/sources.list.d/kubernetes.list

mkdir /run/download
apt update
apt install -y --allow-downgrades \
  kubelet=${VERSION} kubeadm=${VERSION} kubectl=${VERSION} \
  docker.io \
  awscli \
  jq

kubeadm config images pull &

wget https://github.com/mikefarah/yq/releases/download/3.4.1/yq_linux_amd64 -O /run/download/yq
install /run/download/yq /usr/bin

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
