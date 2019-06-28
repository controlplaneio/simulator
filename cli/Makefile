NAME := simulator-cli-build
PACKAGE_NAME := simulator-standalone/cli
BIN_NAME := simulator
GITHUB_ORG = controlplaneio

PKG := github.com/$(GITHUB_ORG)/$(PACKAGE_NAME)

SHELL := /bin/bash
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

GIT_COMMIT := $(shell git -c log.showSignature=false rev-parse HEAD)
GIT_UNTRACKED_CHANGES := $(shell git -c log.showSignature=false \
	status --porcelain)

ifneq ($(GIT_UNTRACKED_CHANGES),)
  GIT_COMMIT := $(GIT_COMMIT)-dirty
endif

# golang buildtime, more at https://github.com/jessfraz/pepper/blob/master/Makefile
# BUG: (rem) this is broken because the
CTIMEVAR=-X $(PKG)/version.GITCOMMIT=$(GITCOMMIT) -X $(PKG)/version.VERSION=$(VERSION)
GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

export NAME BUILD_DATE GIT_SHA
### github.com/controlplaneio/ensure-content.git makefile-header END ###

.PHONY: all
all: help

.PHONY: test
test: test-unit test-acceptance ## run all tests

.PHONY: test-acceptance
test-acceptance: build ## run bats acceptance tests
	@echo "+ $@"
	bash -xc 'cd test && ./bin/bats/bin/bats $(BATS_PARALLEL_JOBS) .'

.PHONY: test-unit
test-unit: build ## run golang unit tests
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
	godocdown pkg/runner > docs/runner.md
	./bin/generate-readme.sh > README.md

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
