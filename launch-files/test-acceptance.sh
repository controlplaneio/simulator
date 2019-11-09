#!/usr/bin/env bash

echo "Running acceptance tests"
mkdir -p ~/.ssh
goss validate
./test/run-tests.tcl

