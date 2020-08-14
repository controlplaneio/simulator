#!/bin/bash

write_keys() {
  set -x
  (
    for USER in ${@:-}; do
      curl --fail --max-time 10 "https://github.com/${USER}.keys" || {
        sleep 2
        curl --fail --max-time 10 "https://github.com/${USER}.keys"
      } || true
    done
  ) >>~/.ssh/authorized_keys

  STATUS=$?

  # dedupe
  awk '!x[$0]++' ~/.ssh/authorized_keys \
    >~/.ssh/authorized_keys.new \
    && mv ~/.ssh/authorized_keys.new ~/.ssh/authorized_keys

  return ${STATUS}
}

