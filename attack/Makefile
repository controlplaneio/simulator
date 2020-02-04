NAME := simulator
CONTAINER_BASE_NAME := simulator-attack
GITHUB_ORG := controlplaneio
DOCKER_HUB_ORG ?= controlplane
VERSION := 0.1-dev

LOCAL_DEV_SIMULATION := $(abspath ../simulation-scripts/scenario/etcd-inverted-wedge)

include ../prelude.mk

.DEFAULT_GOAL := help

.PHONY: all
all: test

# --- DOCKER

.PHONY: docker-run
docker-run: docker-test ## Runs the attack container
	docker run \
		-h attack \
		-v $(LOCAL_DEV_SIMULATION)/tasks.yaml:/tasks.yaml \
		--rm --init -it $(CONTAINER_NAME_LATEST) \
		bash

.PHONY: docker-build
docker-build: ## Builds the attack container
	@docker build -f Dockerfile -t $(CONTAINER_NAME_LATEST) ..

.PHONY: docker-test
docker-test: docker-build ## Run the attack container tests
	@docker run \
		--rm -t $(CONTAINER_NAME_LATEST) \
		goss validate

.PHONY: docker-push
docker-push: docker-test ## pushes the last build docker image
	@echo "+ $@"
	docker push "${CONTAINER_NAME}"
	#docker push "${CONTAINER_NAME_LATEST}"

release: docker-build
ifndef RELEASE_TAG
	$(error RELEASE_TAG is undefined)
endif
	docker tag $(CONTAINER_NAME_LATEST) $(DOCKER_HUB_ORG)/simulator-attack:$(RELEASE_TAG)
	docker push "${CONTAINER_NAME_LATEST}"
	docker push $(DOCKER_HUB_ORG)/simulator-attack:$(RELEASE_TAG)

