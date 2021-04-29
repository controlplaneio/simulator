#!/bin/bash

set -Eeuxo pipefail

IMAGE="control-plane.io/valiant:effort"

main() {
  setup
  create_user_and_k8s_secret
  build_image
  docker_run
}

docker_run() {
  docker run --restart=always \
    -p 5678:5678 \
    "${IMAGE}" \
    -text='☣☠☣ This is a valiant effort! ☣☠☣ Follow the 22 white rabbits ✰ Malodorous regards, Hλ$ħ𝔍Ⱥ¢k // https://securi.fyi ✰'
}

prep_file() {
  cat <<EOF >test
flag_ctf{d241b8126e0adff0}
EOF
}

create_user_and_k8s_secret() {
  adduser hashjack docker
  SSH_KEY_PATH="/home/hashjack/.ssh/id_rsa"
  rm -rf "${SSH_KEY_PATH}"
  su hashjack -c "id; ssh-keygen -f ${SSH_KEY_PATH} -N ''"
  cat /home/hashjack/.ssh/id_rsa.pub >> /home/hashjack/.ssh/authorized_keys
  chown -R hashjack:hashjack /home/hashjack/

  if kubectl get secret my-secret; then
    kubectl delete secret my-secret
  fi
  kubectl create secret generic my-secret \
    --from-file=ssh-privatekey="${SSH_KEY_PATH}"
}

build_image() {
  prep_file
  cat <<EOF | docker build . -f - --tag "${IMAGE}"
FROM hashicorp/http-echo
COPY test /proc/self/cmdline
EOF
}

setup() {
  DIR=$(mktemp -d)
  cd "${DIR}"
  trap cleanup EXIT
}

cleanup() {
  rm -rf "${DIR}"
}

# main &> /tmp/logs
main
