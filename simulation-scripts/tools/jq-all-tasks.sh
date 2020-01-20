#! /bin/bash

for dir in simulation-scripts/scenario/*/ ; do
  echo "${dir}"
  new_yaml=$(yq r -j "${dir}"tasks.yaml | jq '.tasks[].hints |= .+[{"summary": ""}]' | yq r -)
  echo "${new_yaml}" > "${dir}"tasks.yaml
done