ifeq ($(NAME),)
  $(error NAME required, please add "NAME := project-name" to top of Makefile)
else ifeq ($(GITHUB_ORG),)
    $(error GITHUB_ORG required, please add "GITHUB_ORG := controlplaneio" to top of Makefile)
else ifeq ($(DOCKER_HUB_ORG),)
    $(error DOCKER_HUB_ORG required, please add "DOCKER_HUB_ORG := controlplane" to top of Makefile)
endif

DOCKER_REGISTRY_FQDN ?= docker.io
DOCKER_HUB_URL := $(DOCKER_REGISTRY_FQDN)/$(DOCKER_HUB_ORG)/$(NAME)

SHELL := /usr/bin/env bash
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

GIT_MESSAGE := $(shell git -c log.showSignature=false \
	log --max-count=1 --pretty=format:"%H")
GIT_SHA := $(shell git -c log.showSignature=false rev-parse HEAD)
GIT_TAG := $(shell bash -c 'TAG=$$(git -c log.showSignature=false \
	describe --tags --exact-match --abbrev=0 $(GIT_SHA) 2>/dev/null); echo "$${TAG:-dev}"')
GIT_UNTRACKED_CHANGES := $(shell git -c log.showSignature=false \
	status --porcelain)

ifneq ($(GIT_UNTRACKED_CHANGES),)
  GIT_COMMIT := $(GIT_SHA)-dirty
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
