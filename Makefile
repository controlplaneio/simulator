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

SSH_CONFIG_PATH := $(HOME)/.ssh/
KUBE_SIM_TMP := $(HOME)/.kubesim/
SIMULATOR_CONFIG_FILE := $(KUBE_SIM_TMP)/simulator.yaml
HOST := $(shell hostname)

# --- Make
.DEFAULT_GOAL := help

.PHONY: all
all: test

# --- Setup and helpers
.PHONY: setup-dev
setup-dev: ## Initialise simulation tree with git hooks
	@ln -s $(shell pwd)/setup/hooks/pre-commit $(shell pwd)/.git/hooks/pre-commit

.PHONY: devtools
devtools: ## Install devtools
	cd tools/scenario-tools && npm install && npm link

.PHONY: reset
reset: ## Clean up files left over by simulator
	@rm -rf ~/.ssh/cp_simulator_*
	@rm -f -- ~/.kubesim/*

.PHONY: validate-requirements
validate-requirements: ## Verify all requirements are met
	@./scripts/validate-requirements

.PHONY: previous-tag
previous-tag:
	@echo The previously released tag was $$(git describe --abbrev=0 --tags)

#previous-tag and release-tag need to be seperate due to $(eval ...) evaluating before any other commands in the recipe no matter the command ordering.
#Otherwise the prompt is displayed before the prev-tag info is echoed

.PHONY: release-tag
release-tag:
	$(eval RELEASE_TAG := $(shell read -p "Tag to release: " tag; echo $$tag))

.PHONY: gpg-preflight
gpg-preflight:
	@echo Your gpg key for git is $$(git config user.signingkey)

# --- DOCKER
run: validate-requirements docker-build ## Run the simulator - the build stage of the container runs all the cli tests
	@docker run                                                 \
		-h launch                                           \
		-v $(SIMULATOR_AWS_CREDS_PATH):/home/launch/.aws    \
		-v $(SSH_CONFIG_PATH):/home/launch/.ssh             \
		-v $(KUBE_SIM_TMP):/home/launch/.kubesim            \
		--env-file launch-environment                       \
		--rm --init -it $(CONTAINER_NAME_LATEST)

.PHONY: docker-build-no-cache
docker-build-no-cache: ## Builds the launch container
	@mkdir -p ~/.kubesim
	@touch ~/.kubesim/simulator.yaml
	@docker build --no-cache -t $(CONTAINER_NAME_LATEST) .

.PHONY: docker-build
docker-build: ## Builds the launch container
	@mkdir -p ~/.kubesim
	@touch ~/.kubesim/simulator.yaml
	@docker build -t --no-cache $(CONTAINER_NAME_LATEST) .

.PHONY: docker-build
docker-build: ## Builds the launch container
	@mkdir -p ~/.kubesim
	@touch ~/.kubesim/simulator.yaml
	@docker build -t $(CONTAINER_NAME_LATEST) .

.PHONY: docker-test
docker-test: docker-build ## Run the tests
	@docker run                                                 \
		-v "$(SIMULATOR_AWS_CREDS_PATH)":/home/launch/.aws  \
		--env-file launch-environment                       \
		--entrypoint ./acceptance.sh                        \
		--rm -t $(CONTAINER_NAME_LATEST)

# -- SIMULATOR CLI
.PHONY: dep
dep: ## Install dependencies for other targets
	$(GO) get github.com/robertkrimen/godocdown/godocdown 2>&1
	$(GO) mod download 2>&1

.PHONY: build
build: dep ## Run golang build for the CLI program
	@echo "+ $@"
	$(GO) build -a -o ./dist/simulator

.PHONY: is-in-launch-container
is-in-launch-container: ## checks you are running in the launch container
	[ $(HOST) == "launch" ]

.PHONY: test
test:  test-unit test-acceptance ## run all tests except goss tests

.PHONY: test-acceptance
test-acceptance: is-in-launch-container build ## Run tcl acceptance tests for the CLI program
	@echo "+ $@"
	./test/run-tests.tcl

.PHONY: test
test-smoke: build ## Run expect smoke test to check happy path works end-to-end
	@echo "+ $@"
	./test/smoke.expect

.PHONY: test-unit
test-unit: build ## Run golang unit tests for the CLI program
	@echo "NOTE YOU SHOULD RUN THESE WITH make docker-test"
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

.PHONY: release
release: gpg-preflight previous-tag release-tag docker-test docker-build build
	git tag --sign -m $(RELEASE_TAG) $(RELEASE_TAG)
	git push origin $(RELEASE_TAG)
	hub release create -m $(RELEASE_TAG) -a dist/simulator $(RELEASE_TAG)
	docker tag $(CONTAINER_NAME_LATEST) $(DOCKER_HUB_ORG)/simulator:$(RELEASE_TAG)
	docker push $(DOCKER_HUB_ORG)/simulator:$(RELEASE_TAG)
	cd attack && make docker-push


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
