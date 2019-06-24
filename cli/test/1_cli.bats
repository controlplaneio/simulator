#!/usr/bin/env bash

load '_helper'

setup() {
  _global_setup
}

teardown() {
  _global_teardown
}

@test "simulator version - prints version" {
  run ../dist/simulator version
  echo 'version: \n' "${output}" >> "${SIMULATOR_CLI_TEST_OUTPUT}"
  [ "${status}" -eq 0 ]
}

@test "simulator scenario list - prints scenarios" {
  run ../dist/simulator scenario list
  echo 'scenario list: \n' "${output}" >> "${SIMULATOR_CLI_TEST_OUTPUT}"
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator scenario launch - prints the selected scenario" {
  run ../dist/simulator scenario list lazy
  echo 'scenario launch: \n' "${output}" >> "${SIMULATOR_CLI_TEST_OUTPUT}"
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator infra create - prints something" {
  run ../dist/simulator infra create
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator infra status - prints something" {
  run ../dist/simulator infra status
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator infra destroy - prints something" {
  run ../dist/simulator infra destroy
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}
