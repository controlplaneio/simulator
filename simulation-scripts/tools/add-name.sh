#! /bin/bash
# shellcheck disable=SC2086,SC1128,SC2086,SC1009,SC1073,SC1072
# This script will add a name to each tasks file
for dir in simulation-scripts/scenario/*/; do
    echo "${dir}"
    # shellcheck erroneously complains about the jq variable not being double 
    # quoted - this syntax is correct
    new_yaml=$(yq r -j "${dir}"tasks.yaml | jq --arg name "$(basename ${dir})" '. + {"name": $name}' | yq r - --prettyPrint)
    echo "${new_yaml}" >"${dir}"tasks.yaml
done
