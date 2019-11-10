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
  create_ssh_directory
  decode_ssh_key
  ssh_host_keyscan

  unset -f main create_ssh_directory decode_ssh_key ssh_host_keyscan

  exec "${@:-/bin/bash}"
}

create_ssh_directory() {
  if [[ ! -d ~/.ssh ]]; then
    mkdir ~/.ssh
    chmod 700 ~/.ssh
  fi
}

decode_ssh_key() {
  if [[ -n "${BASE64_SSH_KEY:-}" ]]; then
    echo "${BASE64_SSH_KEY:-}" | base64 -d >~/.ssh/id_rsa
    chmod 600 ~/.ssh/id_rsa
    info "Deployed ssh key to $(whoami)'s .ssh/ directory"

    eval "$(ssh-agent)"
    ssh-add ~/.ssh/id_rsa
    info "Added key to ssh-agent"
  fi
}

ssh_host_keyscan() {
  {
    if [[ "${MASTER_IP_ADDRESSES:-}" != "" ]]; then
      echo "${MASTER_IP_ADDRESSES//,/ }"
    fi
    if [[ "${NODE_IP_ADDRESSES:-}" != "" ]]; then
      echo "${NODE_IP_ADDRESSES//,/ }"
    fi
  } | xargs -I{} bash -c \
    'ssh-keygen -f "/root/.ssh/known_hosts" -R {}; ssh-keyscan -H {}' \
    >>~/.ssh/known_hosts 2>/dev/null
}

main "${@}"
