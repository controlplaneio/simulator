#!/usr/bin/env bash

SCENARIODIR="$(dirname -- "$0")/../"

helm template gitea gitea \
    --repo https://dl.gitea.io/charts/ \
    --namespace gitea-system \
    --create-namespace \
    --no-hooks \
    --values "$SCENARIODIR/_scripts/values.yaml" \
    >"$SCENARIODIR/apply/03-gitea-deploy.yaml"

#sed -i -e 's|/data/gitea/conf/app.ini|/config/app.ini|g' "$SCENARIODIR/apply/03-gitea-deploy.yaml"
