ifeq ($(NAME),)
  $(error NAME required, please add "NAME := project-name" to top of Makefile)
else ifeq ($(GITHUB_ORG),)
    $(error GITHUB_ORG required, please add "GITHUB_ORG := controlplaneio" to top of Makefile)
else ifeq ($(DOCKER_HUB_ORG),)
    $(error DOCKER_HUB_ORG required, please add "DOCKER_HUB_ORG := controlplane" to top of Makefile)
endif

DOCKER_REGISTRY_FQDN ?= docker.io

SHELL := /usr/bin/env bash
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

GIT_MESSAGE := $(shell git -c log.showSignature=false \
	log --max-count=1 --pretty=format:"%H")
GIT_SHA := $(shell git rev-parse HEAD)
# GIT_TAG is the tag corresponsing to the curent SHA for HEAD or "dev"
GIT_TAG := $(shell bash -c 'TAG=$$(git -c log.showSignature=false \
	describe --tags --exact-match --abbrev=0 $(GIT_SHA) 2>/dev/null); echo "$${TAG:-dev}"')
GIT_UNTRACKED_CHANGES := $(shell git -c log.showSignature=false \
	status --porcelain)
MOST_RECENT_TAG := $(shell git describe --tags --abbrev=0)

ifneq ($(GIT_UNTRACKED_CHANGES),)
  ifneq ($(GIT_TAG),dev)
    GIT_TAG := $(GIT_TAG)-dirty
  endif
endif

CONTAINER_BASE_NAME ?= $(NAME)
CONTAINER_TAG ?= $(GIT_TAG)
CONTAINER_TAG_LATEST := $(CONTAINER_TAG)
CONTAINER_NAME := $(DOCKER_REGISTRY_FQDN)/$(DOCKER_HUB_ORG)/$(CONTAINER_BASE_NAME):$(CONTAINER_TAG)

# if no untracked changes and tag is not dev, release `latest` tag
ifeq ($(GIT_UNTRACKED_CHANGES),)
  ifneq ($(GIT_TAG),dev)
    CONTAINER_TAG_LATEST = latest
  endif
endif

CONTAINER_NAME_LATEST := $(DOCKER_REGISTRY_FQDN)/$(DOCKER_HUB_ORG)/$(CONTAINER_BASE_NAME):$(CONTAINER_TAG_LATEST)

PKG := github.com/$(GITHUB_ORG)/$(GO_MODULE_NAME)
DATE := $(shell date -I)

GO_LDFLAGS=-ldflags "-w -X $(PKG)/cmd.commit=$(GIT_SHA) -X $(PKG)/cmd.version=$(MOST_RECENT_TAG) -X $(PKG)/cmd.date=$(DATE)"

GO := go

