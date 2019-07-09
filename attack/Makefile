NAME := simulator
CONTAINER_BASE_NAME := simulator-attack
GITHUB_ORG := controlplaneio
DOCKER_HUB_ORG := controlplane
VERSION := 0.1-dev

include ../prelude.mk

.DEFAULT_GOAL := help

.PHONY: all
all: test

# --- DOCKER

.PHONY: docker-run
docker-run: docker-test ## Runs the attack container
	docker run                                  \
		--rm --init -it $(CONTAINER_NAME_LATEST)  \
		bash

.PHONY: docker-build
docker-build: ## Builds the attack container
	@docker build -t $(CONTAINER_NAME_LATEST) .

.PHONY: docker-test
docker-test: docker-build ## Run the attack container tests
	@docker run                        \
		--rm -t $(CONTAINER_NAME_LATEST) \
		goss validate

.PHONY: docker-push
docker-push: docker-test ## pushes the last build docker image
	@echo "+ $@"
	docker push "${CONTAINER_NAME}"
	docker push "${CONTAINER_NAME_LATEST}"

