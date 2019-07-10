[ -n "$TERM" ] && attack-motd

COLOUR_WHITE=$(tput setaf 7 :-"" 2>/dev/null)
COLOUR_RED=$(tput setaf 1 :-"" 2>/dev/null)
COLOUR_RESET=$(tput sgr0 :-"" 2>/dev/null)

BOLD=$(tput bold)
NORMAL=$(tput sgr0)

IFS=',' read -r -a master  <<< "$MASTER_IP_ADDRESSES"
export master
IFS=',' read -r -a node  <<< "$NODE_IP_ADDRESSES"
export node

warning() {
  [ "${*:-}" ] && ERROR="$*" || ERROR="Unknown Warning"
  printf "%s\\n" "${COLOUR_RED}${ERROR}${COLOUR_RESET}"
} 1>&2
readonly -f warning

info() {
  [ "${*:-}" ] && INFO="$*" || INFO="Unknown Info"
  printf "%s\\n" "$(log_message_prefix)${COLOUR_WHITE}${INFO}${COLOUR_RESET}"
} 1>&2
readonly -f info

attack_master() {
  [[ $# = 0 ]] && warning "need a number"

  ssh root@"${master[$1]}"
}
readonly -f attack_master
export -f attack_master

attack_node() {
  [[ $# = 0 ]] && warning "need a number"

  ssh root@"${node[$1]}"
}
readonly -f attack_node
export -f attack_node

welcome() {
  info "You have made in into the permiter of a private kubernetes cluster."
  info "\n"
  info "There are ${BOLD}${#master[@]}${NORMAL} masters 2 "   \
       "and ${BOLD}${#node[@]}${NORMAL} nodes in the cluster"
  info "\n"
  info "We have setup helpers to attack the master(s) and node(s)."
  info "\nTo attack a master:"
  info "  \$ attack_master 0"
  info "\nTo attack a node:"
  info "  \$ attack_node 2"
  info "\nTo see this message again:"
  info "  \$ welcome"
}
readonly -f welcome
export -f welcome
