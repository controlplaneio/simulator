#!/usr/bin/env bash

kubectl delete node k8s-node-1 --force --grace-period=0 --ignore-not-found

# Prep git repos
mkdir -p /tmp/{git,ci}repo
