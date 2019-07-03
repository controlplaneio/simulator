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
  generate_readme
  generate_cli_usage
}
readonly -f main

generate_readme() {
  local -r template="$(cat ./doc-templates/README.template.md)"
  local -r make="$(make -s help-no-color)"

  eval "echo \"${template}\"" > ./README.md
  return 0
}
readonly -f generate_readme

generate_cli_usage() {
  local -r template="$(cat ./doc-templates/cli.template.md)"
  local -r config_help="$(./dist/simulator config help)"
  local -r help="$(./dist/simulator help)"
  local -r infra_help="$(./dist/simulator infra help)"
  local -r scenario_help="$(./dist/simulator scenario help)"
  local -r ssh_help="$(./dist/simulator ssh help)"

  eval "echo \"${template}\"" > ./docs/cli.md
  return 0
}
readonly -f generate_cli_usage

main
