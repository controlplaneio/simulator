#!/usr/bin/env bash

kubectl delete node {{ node2_hostname }} --force --grace-period=0 --ignore-not-found

# Prep git repos
mkdir -p /tmp/{git,ci}repo
