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
  local -r template="$(cat ./README.template.md)"
  local -r make="$(make -s help-no-color)"
  local -r help="$(./dist/simulator help)"
  local -r scenario_help="$(./dist/simulator scenario help)"
  local -r infra_help="$(./dist/simulator infra help)"

  eval "echo \"${template}\""
}

main
