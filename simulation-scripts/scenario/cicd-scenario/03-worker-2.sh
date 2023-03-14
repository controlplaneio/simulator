#!/usr/bin/env bash

set -Eeuxo pipefail

# Early exit for development re-runs
if [[ -f /tmp/.runner ]]; then
    echo "Skipping non-declarative setup"
    exit 0
fi

systemctl disable --now kubelet

apt update && apt install -y --no-install-recommends docker.io

systemctl enable --now docker

# Get Token
# It's fine to hardcode creds here as this script never touches disk on the node.
USER="ctf_admin"
PASS="ahXeehohsoo2suej4tee0ol5xeeteM1w"

export KUBECONFIG=/etc/kubernetes/kubelet.conf
MASTER_IP=$(kubectl get nodes k8s-master-0 -ojsonpath='{.status.addresses}' | jq '.[] | select(.type=="InternalIP") | .address' -r)
BASEURL="http://$MASTER_IP:30080/"

# Gitea may take a minute to come up here; allow 120s
CSRF=$(curl -sS --fail -c cookie.jar "$BASEURL" --retry 12 --retry-connrefused --retry-delay 10 | awk -F' ' '/csrfToken/ {print $2}' | tr -d "',")
curl -sSL -b cookie.jar -c cookie.jar -XPOST "$BASEURL/user/login" -d "user_name=$USER&password=$PASS&_csrf=$CSRF"
TOKEN=$(curl -sSL --fail -b cookie.jar "$BASEURL/admin/runners" | grep -A3 'Registration Token' | awk -F'"' '/value/ {print $4}')

# Disable Kubelet creds
rm -f /var/lib/kubelet/pki/kubelet-client-current.pem
shred -fu /var/lib/kubelet/pki/*

touch /tmp/.runner
docker run --name act-runner -d --restart always -e GITEA_CI_TOKEN="$TOKEN" -e GITEA_CI_URL="$BASEURL" --uts host -v /var/run/docker.sock:/var/run/docker.sock -v /tmp/.runner:/runner/.runner docker.io/controlplaneoffsec/act-runner:latest
docker pull docker.io/library/node:16-bullseye
