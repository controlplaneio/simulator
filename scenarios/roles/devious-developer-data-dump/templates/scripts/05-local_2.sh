
set -x

USER="jcastillo"
PASS="9zvB2cQf2tdC"
ORG="rescue-drop"
DOMAIN="rescue.drop"
REPO="production-image-build"
MASTER_IP="{{ master_ip }}"
NODE_IP="{{ node1_ip }}"


sed -i -e "s/__REGISTRY_IP__/$MASTER_IP/g" /etc/containerd/certs.d/reg.rescue.drop/hosts.toml
systemctl restart containerd
