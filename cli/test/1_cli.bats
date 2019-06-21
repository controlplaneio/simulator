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
  [ "${status}" -eq 0 ]
}

@test "simulator scenario list - prints scenarios" {
  run ../dist/simulator scenario list
  echo "${output}" >&3
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}

@test "simulator scenario launch - prints scenarios" {
  run ../dist/simulator scenario list lazy
  echo "${output}" >&3
  [ "${output}" != "" ]
  [ "${status}" -eq 0 ]
}
