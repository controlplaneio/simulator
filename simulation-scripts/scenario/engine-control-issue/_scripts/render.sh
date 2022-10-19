#!/bin/bash

SCENARIODIR="$(dirname -- "$0")/../"

helm template ingress-nginx-one ingress-nginx \
    --repo https://kubernetes.github.io/ingress-nginx \
    --namespace engine-x-line-one \
    --values "$SCENARIODIR/_scripts/one.yaml" \
    >"$SCENARIODIR/apply/01-deploy-nginx-ingress-one.yaml"

helm template ingress-nginx-two ingress-nginx \
    --repo https://kubernetes.github.io/ingress-nginx \
    --namespace engine-x-line-two \
    --values "$SCENARIODIR/_scripts/two.yaml" \
    >"$SCENARIODIR/apply/01-deploy-nginx-ingress-two.yaml"

helm template ingress-nginx-three ingress-nginx \
    --repo https://kubernetes.github.io/ingress-nginx \
    --namespace engine-x-line-three \
    --values "$SCENARIODIR/_scripts/three.yaml" \
    >"$SCENARIODIR/apply/01-deploy-nginx-ingress-three.yaml"
