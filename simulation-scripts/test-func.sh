test_scenario() {

  local SCENARIO="${1:-}"
  local FOUND_SCENARIO="${SCENARIO}"

  if [[ "${FOUND_SCENARIO:-}" == "" ]]; then
    # if no hash found, but we want to find a test script
    if ! FOUND_SCENARIO=$(find_scenario); then

      if [[ "${SCENARIO:-}" != "" ]]; then
        info "Scenario passed to function, using ${SCENARIO}"
        FOUND_SCENARIO="${SCENARIO}"
      else
        warning "No scenario found"
        return 99
      fi
    fi
  fi

  local TEST_SCRIPT="scenario/${FOUND_SCENARIO:-___empty}/test.sh"

  if [[ -f "${TEST_SCRIPT}" ]]; then
    info "Tests running from ${TEST_SCRIPT}"
    ( source "${TEST_SCRIPT}" )
    return $?
  else
    warning "Test script not found at: ${TEST_SCRIPT}"
    return 99
  fi

  return 1
}

find_scenario() {

  local CHALLENGE_HASH
  local FOUND_SCENARIO=""

  CHALLENGE_HASH=$(echo 'cat /opt/challenge.txt && { echo ; base64 -w0 /opt/challenge.txt; }; echo' | run_ssh "$(get_master)" | tail -n1)

  if [[ "${CHALLENGE_HASH:-}" != "" ]]; then
    warning "SCENARIO hash: ${CHALLENGE_HASH}"
    for CHALLENGE in scenario/**/challenge.txt; do
      THIS_CHALLENGE=$(cat "${CHALLENGE}" | base64 -w0)
      if [[ "${THIS_CHALLENGE}" == "${CHALLENGE_HASH}" ]]; then
        FOUND_SCENARIO=$(basename $(dirname "${CHALLENGE}"))
        info "Installed scenario found: ${FOUND_SCENARIO}"
        break
      fi
    done
  else
    info "No hash found"
  fi

  echo "${FOUND_SCENARIO}"
}

read_prompt() {
  local IS_PROCEED
  echo
  read -n 1 -p "Proceed? [n/q/Y] " IS_PROCEED
  local IS_PROCEED=$(echo "${IS_PROCEED}" | tr '[A-Z]' '[a-z]')
  echo
  [[ "${IS_PROCEED}" = n || "${IS_PROCEED}" = q ]] && return 1
  return 0
}
