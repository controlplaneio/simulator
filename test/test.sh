#!/usr/bin/env bash

#
# Prelude - make bash behave sanely
# http://redsymbol.net/articles/unofficial-bash-strict-mode/
#
set -euo pipefail
IFS=$'\n\t'
# Beware of CDPATH gotchas causing cd not to work correctly when a user has
# set this in their environment
# https://bosker.wordpress.com/2012/02/12/bash-scripters-beware-of-the-cdpath/
unset CDPATH

#
# Main function
#
main() {
  build_launch_container
  run_tests
  exit 0
}

build_launch_container() {
  pushd ..
  docker build -t launch-container .
  popd
}

run_tests() {
  dgoss run \
    -v "${HOME}/.aws/credentials:/app/credentials" \
    -t launch-container || exit 1
}

main
