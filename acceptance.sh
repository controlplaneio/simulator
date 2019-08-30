#!/usr/bin/env bash

echo "Running acceptance tests"
ls -lasp
goss validate
./test/run-tests.tcl

