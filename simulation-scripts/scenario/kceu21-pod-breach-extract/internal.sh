#!/bin/bash

set -Eeuxo pipefail

main() {
  setup

  if [[ -d /var/local/.kube ]]; then
    mv /var/local/.kube /root
  fi

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
