#!/bin/bash

echo "Running acceptance tests"
mkdir -p ~/.ssh
goss validate
./test/run-tests.tcl

