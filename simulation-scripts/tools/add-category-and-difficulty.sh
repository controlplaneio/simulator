#! /bin/bash
# APPEARS TO BROKEN DO NOT USE
# This script will find the category and difficulty of each scenario and add these values to new fields in the tasks.yaml.
# for dir in simulation-scripts/scenario/*/; do
#     echo "${dir}"
#     category=$(echo "${dir}" | cut -f3 -d"/" | cut -f1 -d"-")
#     difficulty=$(grep Difficulty "${dir}"challenge.txt | cut -f2 -d":" | tr -d " ")
#     echo "${category} ${difficulty}"
#     new_yaml=$(yq r -j "${dir}"tasks.yaml | jq --arg CATEGORY "${category}" '. |= .+{"category": $CATEGORY}' | jq --arg DIFFICULTY "${difficulty}" '. |= .+{"difficulty": $DIFFICULTY}' | yq r - | yq r -j - | jq '.' | yq r -)
#     echo "${new_yaml}" >"${dir}"tasks.yaml
# done
