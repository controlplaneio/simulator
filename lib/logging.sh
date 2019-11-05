#!/usr/bin/env bash

COLOUR_RED=$(tput setaf 1 :-"" 2>/dev/null)
COLOUR_GREEN=$(tput setaf 2 :-"" 2>/dev/null)
COLOUR_WHITE=$(tput setaf 7 :-"" 2>/dev/null)
COLOUR_RESET=$(tput sgr0 :-"" 2>/dev/null)

success() {
  [ "${*:-}" ] && RESPONSE="$*" || RESPONSE="Unknown Success"
  printf "%s\n" "$(log_message_prefix)${COLOUR_GREEN}${RESPONSE}${COLOUR_RESET}"
} 1>&2
readonly -f success

info() {
  [ "${*:-}" ] && INFO="$*" || INFO="Unknown Info"
  printf "%s\n" "$(log_message_prefix)${COLOUR_WHITE}${INFO}${COLOUR_RESET}"
} 1>&2
readonly -f info

warning() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Warning"
  printf "%s\n" "$(log_message_prefix)${COLOUR_RED}${ERROR}${COLOUR_RESET}"
} 1>&2
readonly -f warning

error() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Error"
  printf "%s\n" "$(log_message_prefix)${COLOUR_RED}${ERROR}${COLOUR_RESET}"
  exit 3
} 1>&2
readonly -f error

log_message_prefix() {
  TIMESTAMP="[$(date +'%Y-%m-%dT%H:%M:%S%z')]"
  THIS_SCRIPT_SHORT=${THIS_SCRIPT/$DIR/.}
  tput bold 2>/dev/null
  echo -n "${TIMESTAMP} ${THIS_SCRIPT_SHORT}: "
}
readonly -f log_message_prefix


