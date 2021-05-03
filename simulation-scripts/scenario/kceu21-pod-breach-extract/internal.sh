#!/bin/bash

set -Eeuxo pipefail

IMAGE="control-plane.io/valiant:effort"

main() {
  setup

  if [[ ! -d /root/.kube ]]; then
    sleep 30
    ls -lasp /var/local/.kube || true
    if [[ ! -d /root/.kube ]]; then
      sleep 30
      ls -lasp /var/local/.kube || true
      if [[ ! -d /root/.kube ]]; then
        sleep 30
        ls -lasp /var/local/.kube || true
        if [[ ! -d /root/.kube ]]; then
          echo "/root/.kube not found"
          if [[ -d /var/local/.kube ]]; then
            ls -lasp /var/local/.kube
            mv /var/local/.kube /root
          else
            exit 1
          fi
        fi
      fi
    fi
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
