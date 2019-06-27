NAME := simulator-cli-build
PACKAGE_NAME := simulator-standalone/cli
BIN_NAME := simulator
GITHUB_ORG = controlplaneio
DOCKER_HUB_ORG = controlplane

PKG := github.com/$(GITHUB_ORG)/$(PACKAGE_NAME)
DOCKER_REGISTRY_FQDN ?= docker.io
DOCKER_HUB_URL := $(DOCKER_REGISTRY_FQDN)/$(DOCKER_HUB_ORG)/$(NAME)

SHELL := /bin/bash
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

GIT_MESSAGE := $(shell git -c log.showSignature=false \
	log --max-count=1 --pretty=format:"%H")
GIT_SHA := $(shell git -c log.showSignature=false rev-parse HEAD)
GIT_TAG := $(shell bash -c 'TAG=$$(git -c log.showSignature=false \
	describe --tags --exact-match --abbrev=0 $(GIT_SHA) 2>/dev/null); echo "$${TAG:-dev}"')
GIT_UNTRACKED_CHANGES := $(shell git -c log.showSignature=false \
	status --porcelain)

ifneq ($(GIT_UNTRACKED_CHANGES),)
  GIT_COMMIT := $(GIT_COMMIT)-dirty
  ifneq ($(GIT_TAG),dev)
    GIT_TAG := $(GIT_TAG)-dirty
  endif
endif

CONTAINER_TAG ?= $(GIT_TAG)
CONTAINER_TAG_LATEST := $(CONTAINER_TAG)
CONTAINER_NAME := $(DOCKER_REGISTRY_FQDN)/$(NAME):$(CONTAINER_TAG)

# if no untracked changes and tag is not dev, release `latest` tag
ifeq ($(GIT_UNTRACKED_CHANGES),)
  ifneq ($(GIT_TAG),dev)
    CONTAINER_TAG_LATEST = latest
  endif
endif

CONTAINER_NAME_LATEST := $(DOCKER_REGISTRY_FQDN)/$(NAME):$(CONTAINER_TAG_LATEST)

# golang buildtime, more at https://github.com/jessfraz/pepper/blob/master/Makefile
CTIMEVAR=-X $(PKG)/version.GITCOMMIT=$(GITCOMMIT) -X $(PKG)/version.VERSION=$(VERSION)
GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

export NAME DOCKER_HUB_URL BUILD_DATE GIT_MESSAGE GIT_SHA GIT_TAG \
  CONTAINER_TAG CONTAINER_NAME CONTAINER_NAME_LATEST CONTAINER_NAME_TESTING
### github.com/controlplaneio/ensure-content.git makefile-header END ###

.PHONY: all
all: help

.PHONY: test
test: test-unit test-acceptance## unit and local acceptance tests

.PHONY: test-acceptance
test-acceptance: build ## acceptance tests
	@echo "+ $@"
	bash -xc 'cd test && ./bin/bats/bin/bats $(BATS_PARALLEL_JOBS) .'

.PHONY: test-go-fmt
test-go-fmt: ## golang fmt check
	@echo "+ $@"
	gofmt -l -s ./pkg | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi
	gofmt -l -s ./cmd | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

.PHONY: test-unit
test-unit: build ## golang unit tests
	@echo "+ $@"
	go test -race $$(go list ./... | grep -v '/vendor/') -run "$${RUN:-.*}"

.PHONY: test-unit-verbose
test-unit-verbose: build ## golang unit tests (verbose)
	@echo "+ $@"
	go test -race -v $$(go list ./... | grep -v '/vendor/') -run "$${RUN:-.*}"

# ---
.PHONY: dep
dep: ## install dependencies for other targets
	go get github.com/robertkrimen/godocdown/godocdown

.PHONY: build
build: dep ## golang build
	@echo "+ $@"
	go build -a -o ./dist/simulator

.PHONY: doc
doc: dep ## generate markdown documentation for packages
	@echo "+ $@"
	godocdown pkg/scenario > docs/scenario.md
	./bin/generate-readme.sh > README.md

# ---

.PHONY: docker-build
docker-build: ## builds a docker image
	@echo "+ $@"
	docker build --tag "${CONTAINER_NAME}" .
	docker tag "${CONTAINER_NAME}" "${CONTAINER_NAME_LATEST}"
	@echo "Successfully tagged ${CONTAINER_NAME} as ${CONTAINER_NAME_LATEST}"
	docker run -v "${PWD}:/go/src/github.com/controlplaneio/simulator-standalone/cli" -t "${CONTAINER_NAME_LATEST}" make test

# ---

.PHONY: help
help: ## parse jobs and descriptions from this Makefile
	set -x;
	@grep -E '^[ a-zA-Z0-9_-]+:([^=]|$$)' $(MAKEFILE_LIST) \
    | grep -Ev '^help\b[[:space:]]*:' \
    | sort \
    | awk 'BEGIN {FS = ":.*?##"}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: help
help-no-color:
	@set -x;
	@grep -E '^[ a-zA-Z0-9_-]+:([^=]|$$)' $(MAKEFILE_LIST) \
    | grep -Ev '^help\b[[:space:]]*:' \
    | sort \
    | awk 'BEGIN {FS = ":.*?##"}; {printf "%-20s %s\n", $$1, $$2}'
