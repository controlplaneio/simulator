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
# colours
#
export CLICOLOR=1
export TERM=xterm-color
COLOUR_RED=$(tput setaf 1 :-"" 2>/dev/null)
COLOUR_GREEN=$(tput setaf 2 :-"" 2>/dev/null)
COLOUR_WHITE=$(tput setaf 7 :-"" 2>/dev/null)
COLOUR_RESET=$(tput sgr0 :-"" 2>/dev/null)

#
# Logging
#
success() {
  [ "${*:-}" ] && RESPONSE="$*" || RESPONSE="Unknown Success"
  printf "%s\\n" "$(log_message_prefix)${COLOUR_GREEN}${RESPONSE}${COLOUR_RESET}"
} 1>&2
readonly -f success

info() {
  [ "${*:-}" ] && INFO="$*" || INFO="Unknown Info"
  printf "%s\\n" "$(log_message_prefix)${COLOUR_WHITE}${INFO}${COLOUR_RESET}"
} 1>&2
readonly -f info

warning() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Warning"
  printf "%s\\n" "$(log_message_prefix)${COLOUR_RED}${ERROR}${COLOUR_RESET}"
} 1>&2
readonly -f warning

error() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Error"
  printf "%s\\n" "$(log_message_prefix)${COLOUR_RED}${ERROR}${COLOUR_RESET}"
  exit 3
} 1>&2
readonly -f error

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
