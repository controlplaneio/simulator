
#! /bin/bash
# This script will add a name to each tasks file
for dir in simulation-scripts/scenario/*/; do
    echo "${dir}"
    new_yaml=$(yq r -j "${dir}"tasks.yaml | jq --arg name "$(basename $dir)" '. + {"name": $name}' | yq r - --prettyPrint)
    echo "${new_yaml}" >"${dir}"tasks.yaml
done
