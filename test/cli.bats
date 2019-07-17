#!/usr/bin/env bash

load '_helper'

setup() {
  _global_setup
}

teardown() {
  _global_teardown
}

@test "simulator version - prints version" {
  run ${BIN_UNDER_TEST} version
  echo 'version: \n' "${output}" >> "${SIMULATOR_CLI_TEST_OUTPUT}"
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator scenario list - prints scenarios" {
  run ${BIN_UNDER_TEST} scenario list
  echo 'scenario list: \n' "${output}" >> "${SIMULATOR_CLI_TEST_OUTPUT}"
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator scenario launch - prints the selected scenario" {
  run ${BIN_UNDER_TEST} scenario list lazy
  echo 'scenario launch: \n' "${output}" >> "${SIMULATOR_CLI_TEST_OUTPUT}"
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator infra create - prints something" {
  run ${BIN_UNDER_TEST} infra create
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator infra status - prints something" {
  run ${BIN_UNDER_TEST} infra status
  echo '---\ninfra status: \n' "${output}" >> "${SIMULATOR_CLI_TEST_OUTPUT}"
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator infra destroy - prints something" {
  run ${BIN_UNDER_TEST} infra destroy
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator get <key> - prints the key value pair" {
  run ${BIN_UNDER_TEST} config get loglevel
  echo '---\nconfig get: \n' "${output}" >> "${SIMULATOR_CLI_TEST_OUTPUT}"
  [ "${output}" != "logelevel = info\n" ]
  [ "${status}" -eq 0 ]
}

@test "simulator completion - prints something" {
  run ${BIN_UNDER_TEST} completion
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}
