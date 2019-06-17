#!/usr/bin/env bash

load '_helper'

setup() {
  _global_setup
}

teardown() {
  _global_teardown
}

@test "prints version" {
  run _app version
  [ "${status}" -eq 0 ]
}
