#! /bin/bash
# This script will add a name to each tasks file
for dir in simulation-scripts/scenario/*/; do
    echo "${dir}"
    # shellcheck erroneously complains about the jq variable not being double 
    # quoted - this syntax is correct
    # shellcheck disable=SC2086
    new_yaml=$(yq r -j "${dir}"tasks.yaml | jq --arg name "$(basename ${dir})" '. + {"name": $name}' | yq r - --prettyPrint)
    echo "${new_yaml}" >"${dir}"tasks.yaml
done
