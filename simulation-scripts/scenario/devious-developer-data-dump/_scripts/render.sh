#!/usr/bin/env bash

SCENARIODIR="$(dirname -- "$0")/../"

helm template gitea gitea \
    --repo https://dl.gitea.io/charts/ \
    --namespace gitea-system \
    --create-namespace \
    --no-hooks \
    --values "$SCENARIODIR/_scripts/values.yaml" \
    >"$SCENARIODIR/apply/03-gitea-deploy.yaml"

helm template ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace \
    --values "$SCENARIODIR/_scripts/values-ingress.yaml" \
    >"$SCENARIODIR/apply/05-ingress-deploy.yaml"
