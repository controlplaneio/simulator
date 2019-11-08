#!/usr/bin/env bash

set -euxo pipefail

main() {

  trap show_exit_warning EXIT
  trap show_exit_warning SIGTERM

  fix_aws_environment_variables
  write_bash_aliases
  write_inputrc

  exec bash
}

fix_aws_environment_variables() {
  export AWS_REGION="${AWS_REGION:-${AWS_DEFAULT_REGION:-}}"
}

write_bash_aliases() {
  cat >>/home/launch/.bash_aliases <<EOF
alias ll='ls -lasp'
EOF
}

write_inputrc() {
  cat >>/home/launch/.inputrc <<EOF
# Make Tab autocomplete regardless of filename case
set completion-ignore-case on

# List all matches in case multiple possible completions are possible
set show-all-if-ambiguous on
EOF
}

draw_box() {
  local ARGUMENTS=("${@}") LINE MAX_WIDTH
  for THIS_ARGUMENT in "${ARGUMENTS[@]}"; do
    ((MAX_WIDTH < ${#THIS_ARGUMENT})) && {
      LINE="${THIS_ARGUMENT}"
      MAX_WIDTH="${#THIS_ARGUMENT}"
    }
  done
  tput bold
  tput setaf 3
  echo "    -${LINE//?/-}-
   | ${LINE//?/ } |"
  for THIS_ARGUMENT in "${ARGUMENTS[@]}"; do
    printf '   | %s%*s%s |\n' "$(tput setaf 1)" "-${MAX_WIDTH}" "${THIS_ARGUMENT}" "$(tput setaf 3)"
  done
  echo "   | ${LINE//?/ } |
    -${LINE//?/-}-"
  tput sgr 0
}

show_exit_warning() {
  echo "$(tput setaf 2)
||====================================================================||
||//$\\\\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\//$\\\\||
||(100)==================| FEDERAL RESERVE NOTE |================(100)||
||\\\\$//        ~         '------========--------'                \\\\$//||
||<< /        /$\              // ____ \\\\                         \\ >>||
||>>|  12    //L\\\\            // ///..) \\\\         L38036133B   12 |<<||
||<<|        \\\\ //           || <||  >\  ||                        |>>||
||>>|         \\\$/            ||  \$\$ --/  ||        One Hundred     |<<||
||<<|      L38036133B        *\\\\  |\_/  //* series                 |>>||
||>>|  12                     *\\\\/___\_//*   1989                  |<<||
||<<\      Treasurer     ______/Franklin\________     Secretary 12 />>||
||//$\                 ~|UNITED STATES OF AMERICA|~               /$\\\\||
||(100)===================  ONE HUNDRED DOLLARS =================(100)||
||\\\\$//\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\/\\\\$//||
||====================================================================||
$(tput sgr0)"

  draw_box "   If you created any infrastructure and did not destroy it " \
    "   you will be accruing charges in your AWS account"
}

main
