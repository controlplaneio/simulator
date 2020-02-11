#! /bin/bash
# This script will add a penalty value to each hint object.
# In order to update all penalty values, simply re-run the script with the new value replacing '10'. This may cause some cosmetic key reordering in the final yaml.
for dir in simulation-scripts/scenario/*/; do
    echo "${dir}"
    category=$(echo "${dir}" | cut -f3 -d"/" | cut -f1 -d"-")
    difficulty=$(cat "${dir}"challenge.txt | grep Difficulty | cut -f2 -d":" | tr -d " ")
    echo "${category} ${difficulty}"
    new_yaml=$(yq r -j "${dir}"tasks.yaml | jq --arg CATEGORY "${category}" '. |= .+{"category": $CATEGORY}' | jq --arg DIFFICULTY "${difficulty}" '. |= .+{"difficulty": $DIFFICULTY}' | yq r - | yq r -j - | jq '.' | yq r -)
    echo "${new_yaml}" >"${dir}"tasks.yaml
done
