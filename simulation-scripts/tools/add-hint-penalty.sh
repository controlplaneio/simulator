#! /bin/bash
# This script will add a penalty value to each hint object.
# In order to update all penalty values, simply re-run the script with the new value replacing '10'. This may cause some cosmetic key reordering in the final yaml.
for dir in simulation-scripts/scenario/*/; do
    echo "${dir}"
    new_yaml=$(yq r -j "${dir}"tasks.yaml | jq '.tasks[].hints[] |= .+{"penalty": 10}' | yq r -)
    echo "${new_yaml}" >"${dir}"tasks.yaml
done
