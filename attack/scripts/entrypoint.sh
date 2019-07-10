#!/usr/bin/env bash

# shellcheck disable=SC2155


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
# Colours
#
COLOUR_WHITE=$(tput setaf 7 :-"" 2>/dev/null)
COLOUR_RESET=$(tput sgr0 :-"" 2>/dev/null)

#
# Logging
#

info() {
  [ "${*:-}" ] && INFO="$*" || INFO="Unknown Info"
  printf "%s\\n" "${COLOUR_WHITE}${INFO}${COLOUR_RESET}"
} 1>&2
readonly -f info

#
# Main function
#
main() {
  decode_ssh_key
  bash
}
readonly -f main

decode_ssh_key() {
  if [[ -n "${BASE64_SSH_KEY}" ]]; then
    if [ ! -d ~/.ssh ]; then
      mkdir ~/.ssh
      chmod 700 ~/.ssh
    fi
    echo "${BASE64_SSH_KEY}" | base64 -d >~/.ssh/id_rsa
    chmod 600 ~/.ssh/id_rsa
    info "Deployed ssh key to $(whoami)'s .ssh/ directory"

    eval "$(ssh-agent)"
    ssh-add ~/.ssh/id_rsa
    info "Added key to ssh-agent"
  fi

}
readonly -f decode_ssh_key

main
