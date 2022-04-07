#!/bin/bash

#
# Prelude - make bash behave sanely
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
#
set -Eeuo pipefail
IFS=$'\n\t'
# Beware of CDPATH gotchas causing cd not to work correctly when a user has
# set this in their environment
# https://bosker.wordpress.com/2012/02/12/bash-scripters-beware-of-the-cdpath/
unset CDPATH

readonly GOSS_VERSION=v0.3.7

#
# Main function
#
main() {
  echo "Starting at $(date)"
  install_goss
  wait_for_cloud_init
  goss validate -f documentation \
    --sleep 10s \
    --retry-timeout 240s
}
readonly -f main

install_goss() {
  echo "Installing goss ${GOSS_VERSION}"
  curl -sL https://github.com/aelsabbahy/goss/releases/download/${GOSS_VERSION}/goss-linux-amd64 -o /usr/local/bin/goss
  chmod +rx /usr/local/bin/goss
}
readonly -f install_goss

wait_for_cloud_init() {
  echo "Waiting for /var/log/cloud-init.log to complete..."
  while ! grep -q "finish: modules-final: SUCCESS: running modules for final" /var/log/cloud-init.log; do
    echo "Waiting 5s for cloud-init to finish $(tail /var/log/cloud-init.log)"
    sleep 5
    if [[ "${SECONDS}" -gt 300 ]]; then
      echo "WARNING: Waited 300s for cloud init, exiting loop"
      break
    fi
  done
}
readonly -f wait_for_cloud_init

main |& tee -a /tmp/cloud-init.log
