# --- Project configuration
NAME := simulator
GITHUB_ORG := controlplaneio
DOCKER_HUB_ORG := controlplane
VERSION := 0.7-pre

# --- Boilerplate
include prelude.mk

# --- Mount paths
ifeq ($(AWS_SHARED_CREDENTIALS_FILE),)
	SIMULATOR_AWS_CREDS_PATH := $(HOME)/.aws/
else
	SIMULATOR_AWS_CREDS_PATH := $(shell dirname $(AWS_SHARED_CREDENTIALS_FILE))
endif

SIMULATOR_CONFIG_FILE := $(PWD)/simulator.yaml
SSH_CONFIG_PATH := $(HOME)/.ssh/
KUBE_SIM_TMP := $(HOME)/.kubesim/

# --- Make
.DEFAULT_GOAL := help

.PHONY: all
all: test

# --- Setup and helpers
.PHONY: setup-dev
setup-dev: ## Initialise simulation tree with git hooks
	@ln -s $(shell pwd)/setup/hooks/pre-commit $(shell pwd)/.git/hooks/pre-commit

.PHONY: reset
reset: ## Clean up files left over by simulator
	@rm -rf ~/.ssh/cp_simulator_*
	@git checkout simulator.yaml

.PHONY: validate-requirements
validate-requirements: ## Verify all requirements are met
	@./scripts/validate-requirements

# --- DOCKER
run: validate-requirements docker-build ## Run the simulator - the build stage of the container runs all the cli tests
	@docker run                                                         \
		-h launch                                                         \
		-v $(SIMULATOR_CONFIG_FILE):/app/simulator.yaml                   \
		-v $(SIMULATOR_AWS_CREDS_PATH):/home/launch/.aws                  \
		-v $(SSH_CONFIG_PATH):/home/launch/.ssh                           \
		-v $(KUBE_SIM_TMP):/home/launch/.kubesim                          \
		--env-file launch-environment                                     \
		--rm --init -it $(CONTAINER_NAME_LATEST)

.PHONY: docker-build
docker-build: ## Builds the launch container
	@docker build -t $(CONTAINER_NAME_LATEST) .

.PHONY: docker-test
docker-test: docker-build ## Run the tests
	@docker run                                                         \
		-v "$(SIMULATOR_AWS_CREDS_PATH)":/home/launch/.aws                \
		-v $(SSH_CONFIG_PATH):/home/launch/.ssh                           \
		--env-file launch-environment                                     \
		--entrypoint goss                                                 \
		--rm -t $(CONTAINER_NAME_LATEST)                                  \
		validate

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
test-acceptance: build ## Run tcl acceptance tests for the CLI program
	@echo "+ $@"
	./test/run-tests.tcl

.PHONY: test
test-smoke: build ## Run expect smoke test to check happy path works end-to-end
	@echo "+ $@"
	./test/smoke.expect

.PHONY: test-unit
test-unit: build ## Run golang unit tests for the CLI program
	@echo "+ $@"
	$(GO) test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-cleanup
test-cleanup: ## cleans up automated test artefacts if e.g. you ctrl-c abort a test run
	@aws s3 rb s3://controlplane-simulator-state-automated-test || true
	@truncate -s 0 simulator-automated-test.yaml || true

.PHONY: coverage
test-coverage:  ## Run golang unit tests with coverage and opens a browser with the results
	@echo "" > count.out
	$(GO) test -covermode=count -coverprofile=count.out ./...
	$(GO) tool cover -html=count.out

.PHONY: docs
docs: dep ## Generate documentation
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
