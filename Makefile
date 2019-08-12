NAME := simulator
GITHUB_ORG := controlplaneio
DOCKER_HUB_ORG := controlplane
VERSION := 0.7-pre

include prelude.mk

# --- AWS

ifeq ($(AWS_SHARED_CREDENTIALS_FILE),)
	SIMULATOR_AWS_CREDS_PATH := $(HOME)/.aws/
else
	SIMULATOR_AWS_CREDS_PATH := $(shell dirname $(AWS_SHARED_CREDENTIALS_FILE))
endif

SIMULATOR_CONFIG_FILE := $(PWD)/simulator.yaml

# --- AWS

SSH_CONFIG_PATH := $(HOME)/.ssh/

.DEFAULT_GOAL := help

.PHONY: all
all: test

# --- DOCKER
run: validate-requirements reset docker-build ## Run the simulator - the build stage of the container runs all the cli tests

	docker run                                                          \
		-h launch                                                         \
		-v $(SIMULATOR_CONFIG_FILE):/app/simulator.yaml                   \
		-v $(SIMULATOR_AWS_CREDS_PATH):/home/launch/.aws                  \
		-v $(SSH_CONFIG_PATH):/home/launch/.ssh                           \
		-e AWS_ACCESS_KEY_ID                                              \
		-e AWS_REGION                                                     \
		-e AWS_DEFAULT_REGION                                             \
		-e AWS_PROFILE                                                    \
		-e AWS_DEFAULT_PROFILE                                            \
		-e AWS_SECRET_ACCESS_KEY                                          \
		--rm --init -it $(CONTAINER_NAME_LATEST)                          \
		bash

.PHONY: docker-build
docker-build: ## Builds the launch container
	@docker build -t $(CONTAINER_NAME_LATEST) .

.PHONY: docker-test
docker-test: docker-build ## Run the tests
	@docker run                                                         \
		-v $(SIMULATOR_AWS_CREDS_PATH):/home/launch/.aws                  \
		--rm -t $(CONTAINER_NAME_LATEST)                                  \
		goss validate

# --- Setup environment

.PHONY: setup-dev
setup-dev: ## Initialise simulation tree with git hooks
	@ln -s $(shell pwd)/setup/hooks/pre-commit $(shell pwd)/.git/hooks/pre-commit

# -- Reset environment

.PHONY: reset
reset: ## Clean up files left over by simulator
	@rm -rf ~/.ssh/cp_simulator_*
	@git checkout simulator.yaml

# -- Check build requirements are met

.PHONY: validate-requirements
validate-requirements: ## Verify all requirements are met
	@$(shell pwd)/scripts/validate_requirements

# -- SIMULATOR CLI

.PHONY: dep
dep: ## Install dependencies for other targets
	$(GO) get github.com/robertkrimen/godocdown/godocdown 2>&1
	$(GO) mod download 2>&1

.PHONY: build
build: dep ## Run golang build for the CLI program
	@echo "+ $@"
	$(GO) build -a -o ./dist/simulator

.PHONY: test
test: test-unit test-acceptance ## run all tests except goss tests

.PHONY: test-acceptance
test-acceptance: build ## Run bats acceptance tests for the CLI program
	@echo "+ $@"
	./test/commands.test

.PHONY: test
test-smoke: build ## Run bats acceptance tests for the CLI program
	@echo "+ $@"
	./test/smoke.test

.PHONY: test-unit
test-unit: build ## Run golang unit tests for the CLI program
	@echo "+ $@"
	$(GO) test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: coverage
coverage:  ## Run golang unit tests with coverage and opens a browser with the results
	@echo "" > count.out
	$(GO) test -covermode=count -coverprofile=count.out ./...
	$(GO) tool cover -html=count.out

.PHONY: doc
doc: dep ## Generate documentation
	@echo "+ $@"
	godocdown pkg/scenario > docs/api/scenario.md
	godocdown pkg/simulator > docs/api/simulator.md
	godocdown pkg/util > docs/api/util.md
	./scripts/generate-docs-from-templates.sh
	./scripts/tf-auto-doc ./terraform

# --- MAKEFILE HELP

.PHONY: help
help: ## parse jobs and descriptions from this Makefile
	@set -x;
	@grep -E '^[ a-zA-Z0-9_-]+:([^=]|$$)' Makefile \
    | grep -Ev '^(help|all|help-no-color)\b[[:space:]]*:' \
    | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: help
help-no-color:
	@set -x;
	@grep -E '^[ a-zA-Z0-9_-]+:([^=]|$$)' Makefile \
    | grep -Ev '^(help|all|help-no-color)\b[[:space:]]*:' \
    | awk 'BEGIN {FS = ":.*?##"}; {printf "%-20s %s\n", $$1, $$2}'
