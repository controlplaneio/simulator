# --- Project configuration
NAME := simulator
GITHUB_ORG := kubernetes-simulator
DOCKER_HUB_ORG ?= controlplane
GO_MODULE_NAME := simulator

# --- Boilerplate
include prelude.mk

# --- Mount paths
ifeq ($(AWS_SHARED_CREDENTIALS_FILE),)
	SIMULATOR_AWS_CREDS_PATH := $(HOME)/.aws/
else
	SIMULATOR_AWS_CREDS_PATH := $(shell dirname $(AWS_SHARED_CREDENTIALS_FILE))
endif

SSH_CONFIG_PATH := $(HOME)/.kubesim/
KUBE_SIM_TMP := $(HOME)/.kubesim/
SIMULATOR_CONFIG_FILE := $(KUBE_SIM_TMP)/simulator.yaml
HOST := $(shell hostname)
TOOLS_DIR := tools/scenario-tools
SIMULATOR_TFVAR_DIR := $(KUBE_SIM_TMP)/settings
COMMAND :=

export GOSUMDB=off

# --- Make
.DEFAULT_GOAL := help

.PHONY: all
all: test

# --- Setup and helpers
.PHONY: setup-dev
setup-dev: ## Initialise simulation tree with git hooks
	@ln -s $(shell pwd)/setup/hooks/pre-commit $(shell pwd)/.git/hooks/pre-commit

.PHONY: devtools-deps
devtools-deps: ## Install devtools dependencies
	cd $(TOOLS_DIR) && npm install

.PHONY: devtools
devtools: devtools-deps ## Install devtools
	cd $(TOOLS_DIR) && npm install && npm link
	@echo "`scenario` should now be on your PATH"

.PHONY: test-devtools
test-devtools: ## Run all the unit tests for the devtools
	cd $(TOOLS_DIR) && npm test

.PHONY: reset
reset: ## Clean up files left over by simulator
	@rm -f -- ~/.kubesim/*

.PHONY: validate-reqs
validate-reqs: ## Verify all requirements are met
	@./scripts/validate-requirements

.PHONY: previous-tag
previous-tag:
	@echo The previously released tag was $$(git describe --abbrev=0 --tags)

#previous-tag and release-tag need to be seperate due to $(eval ...) evaluating before any other commands in the recipe no matter the command ordering.
#Otherwise the prompt is displayed before the prev-tag info is echoed

.PHONY: release-tag
release-tag:
	$(eval RELEASE_TAG := $(shell read -p "Tag to release: " tag; echo $$tag))

.PHONY: release-msg
release-msg:
	$(eval TMP_FILE := $(shell mktemp /tmp/sim-release-msg.XXXXX))
	@echo "$(RELEASE_TAG)" > $(TMP_FILE)
	@echo "" >> $(TMP_FILE)
	@echo "$$(git log --merges --format=%b $$(git describe --tags --abbrev=0)..HEAD)" >> $(TMP_FILE)

.PHONY: gpg-preflight
gpg-preflight:
	@echo Your gpg key for git is $$(git config user.signingkey)

# --- DOCKER
run: validate-reqs docker-build docker-run ## Run the simulator - the build stage of the container runs all the cli tests

docker-run: validate-reqs ## Run the simulator container
	@docker run                                                             \
		-h launch                                                       \
		-v $(SIMULATOR_AWS_CREDS_PATH):/home/launch/.aws                \
		-v $(KUBE_SIM_TMP):/home/launch/.kubesim                        \
		-v "$(shell pwd)/terraform":/app/terraform                      \
		-v "$(shell pwd)/simulation-scripts":/app/simulation-scripts:ro \
		--env-file launch-environment                                   \
		--rm --init -it $(CONTAINER_NAME_LATEST) $(COMMAND)

.PHONY: docker-build-nocache
docker-build-nocache: ## Builds the launch container without the Docker cache
	@mkdir -p ~/.kubesim
	@touch ~/.kubesim/simulator.yaml
	@mkdir -p $(SIMULATOR_TFVAR_DIR)
	@touch $(SIMULATOR_TFVAR_DIR)/bastion.tfvars
	@docker build --no-cache -t $(CONTAINER_NAME_LATEST) .

.PHONY: docker-build
docker-build: ## Builds the launch container
	@mkdir -p ~/.kubesim
	@touch ~/.kubesim/simulator.yaml
	@mkdir -p $(SIMULATOR_TFVAR_DIR)
	@touch $(SIMULATOR_TFVAR_DIR)/bastion.tfvars
	@docker build -t $(CONTAINER_NAME_LATEST) .

.PHONY: docker-test
docker-test: validate-reqs docker-build ## Run the tests
	@export AWS_DEFAULT_REGION="testing propagation to AWS_REGION var"; \
	docker run                                                          \
		-v "$(SIMULATOR_AWS_CREDS_PATH)":/home/launch/.aws          \
		--env-file launch-environment                               \
		--rm -t $(CONTAINER_NAME_LATEST)                            \
		/app/test-acceptance.sh

	cd attack && make docker-test

# -- SIMULATOR CLI
dep: go.mod ## Install dependencies for other targets
	mkdir -p ~/go/bin
	$(GO) mod download
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ~/go/bin v1.23.8

.PHONY: static-analysis
static-analysis: dep ## Runs static analysis tools over golang code for known problem
	golangci-lint run

.PHONY: build
build: static-analysis ## Run golang build for the CLI program
	@echo "+ $@"
	$(GO) build ${GO_LDFLAGS} -a -o ./dist/simulator

.PHONY: is-in-launch
is-in-launch: ## Checks you are running in the launch container
	[ $(HOST) == "launch" ]

.PHONY: test
test:  test-unit test-acceptance ## Run all tests except goss tests

.PHONY: test-acceptance
test-acceptance: is-in-launch build ## Run tcl acceptance tests for the CLI program
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
test-cleanup: ## Cleans up automated test artefacts if e.g. you ctrl-c abort a test run
	@aws s3 rb s3://controlplane-simulator-state-automated-test || true
	@truncate -s 0 simulator-automated-test.yaml || true

.PHONY: coverage
test-coverage:  ## Run golang unit tests with coverage and opens a browser with the results
	@echo "" > count.out
	$(GO) test -covermode=count -coverprofile=count.out ./...
	$(GO) tool cover -html=count.out

.PHONY: docs
docs: ## Generate documentation
	@echo "+ $@"
	./scripts/tf-auto-doc ./terraform

.PHONY: release
release: validate-reqs gpg-preflight previous-tag release-tag docker-test docker-build build release-msg ## Docker container and binary release automation for simulator
	git tag --sign -m $(RELEASE_TAG) $(RELEASE_TAG)
	git push origin $(RELEASE_TAG)
	hub release create -F $(TMP_FILE) -a dist/simulator -a kubesim $(RELEASE_TAG)
	rm -rf $(TMP_FILE)
	docker tag $(CONTAINER_NAME_LATEST) $(DOCKER_HUB_ORG)/simulator:$(RELEASE_TAG)
	docker push $(DOCKER_HUB_ORG)/simulator:$(RELEASE_TAG)
	docker tag $(CONTAINER_NAME_LATEST) $(DOCKER_HUB_ORG)/simulator:latest
	docker push $(DOCKER_HUB_ORG)/simulator:latest
	cd attack && RELEASE_TAG=$(RELEASE_TAG) make release

# --- PACKER BUILD
.PHONY: packer
packer: ## parse jobs and descriptions from this Makefile
	@set -x; cd packer/k8s-base/images/ && packer build image.pkr.hcl

# --- MAKEFILE HELP
.PHONY: help
help: ## parse jobs and descriptions from this Makefile
	@set -x;
	@grep -E '^[ a-zA-Z0-9_-]+:([^=]|$$)' Makefile \
    | grep -Ev '^(help|all|help-no-colour|previous-tag|release-tag|gpg-preflight)\b[[:space:]]*:' \
    | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: help
help-no-colour:
	@set -x;
	@grep -E '^[ a-zA-Z0-9_-]+:([^=]|$$)' Makefile \
    | grep -Ev '^(help|all|help-no-colour|previous-tag|release-tag|gpg-preflight)\b[[:space:]]*:' \
    | awk 'BEGIN {FS = ":.*?##"}; {printf "%-20s %s\n", $$1, $$2}'
