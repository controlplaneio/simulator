.DEFAULT_GOAL := help

SHELL := /usr/bin/env bash

.PHONY: all
all: test

.PHONY: check
check: ## Check required system packages are installed
	@command -v "npm" > /dev/null 2>&1 || echo >2 "Couldn't find npm - please install nodejs"

.PHONY: deps
deps: check ## Install dependencies
	@npm install

.PHONY: test
test: deps ## Run the feature tests
	@cucumber-js

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

