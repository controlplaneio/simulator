#!/bin/bash

set -Eeuxo pipefail

IMAGE="control-plane.io/valiant:effort"

main() {
  setup

  rm -rf /root/.kube
  mv /var/local/.kube /root/

  # double bind mount secretzy
  SECRET_DIR=$(mktemp -d)
  mount -o bind "${SECRET_DIR}" /root/.kube
  SECRET_DIR=$(mktemp -d)
  mount -o bind "${SECRET_DIR}" /root/.kube
}


setup() {
  DIR=$(mktemp -d)
  cd "${DIR}"
  trap cleanup EXIT
}

cleanup() {
  rm -rf "${DIR}"
}

#main | tee /tmp/logs
main
