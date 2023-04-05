#!/usr/bin/env bash

#set -Eeuo pipefail
set -x

PR_ID=${1}
TOKEN=${2}

## Start review
read -r -d '' DATA <<EOF
{
  "body": "starting review...",
  "comments": [],
  "commit_id": "string",
  "event": "PENDING"
}
EOF
REVIEW_ID=$(curl -X 'POST' \
  "http://localhost:8080/api/v1/repos/rescue-drop/production-image-build/pulls/$PR_ID/reviews?access_token=$TOKEN" \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d "$DATA" \
  | jq .id -r )


## Finish review as approved
curl -X 'POST' \
  "http://localhost:8080/api/v1/repos/rescue-drop/production-image-build/pulls/$PR_ID/reviews/$REVIEW_ID?access_token=$TOKEN" \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"body": "approved","event": "APPROVED"}'

## Merge
COMMIT=$(curl -X 'GET' \
  "http://localhost:8080/api/v1/repos/rescue-drop/production-image-build/pulls/$PR_ID/commits?access_token=$TOKEN" \
  -H 'accept: application/json' \
  | jq .[0].sha -r)

read -r -d '' DATA <<EOF
{
  "Do": "merge",
  "MergeCommitID": "1",
  "MergeMessageField": "string",
  "MergeTitleField": "merge pls",
  "delete_branch_after_merge": true,
  "force_merge": true,
  "head_commit_id": "$COMMIT",
  "merge_when_checks_succeed": false
}
EOF
curl -X 'POST' \
  "http://localhost:8080/api/v1/repos/rescue-drop/production-image-build/pulls/$PR_ID/merge?access_token=$TOKEN" \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d "$DATA"
