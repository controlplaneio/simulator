NAME:=simulator
DOCKER_REPOSITORY:=controlplane
LAUNCH_DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)-launch
VERSION:=0.1-dev
SIMULATOR_CREDS_PATH:=$(HOME)/.aws/credentials

SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

.PHONY: all
all: test

.PHONY: run
run: build ## Runs the simulator
	docker run -v $(SIMULATOR_AWS_CREDS_PATH):/app/credentials --rm -it $(LAUNCH_DOCKER_IMAGE_NAME):$(VERSION) bash

.PHONY: run
exec: build ## Run a command in the launch container - CMD=<...> make exec
	docker run -v $(SIMULATOR_AWS_CREDS_PATH):/app/credentials --rm -it $(LAUNCH_DOCKER_IMAGE_NAME):$(VERSION) $(CMD)

.PHONY: build
build: ## Builds the launch container
	@docker build -t $(LAUNCH_DOCKER_IMAGE_NAME):$(VERSION) .

.PHONY: test
test: build ## Run the tests
	@docker run -v $(SIMULATOR_AWS_CREDS_PATH):/app/credentials --rm -t $(LAUNCH_DOCKER_IMAGE_NAME):$(VERSION) goss validate

.PHONY: infra-init
infra-init:
	@pushd terraform; terraform init; popd

.PHONY: infra-checkvars
infra-checkvars:
	@test -f terraform/settings/bastion.tfvars || \
		(echo Please create terraform/settings/bastion.tfvars && exit 1)

.PHONY: infra-plan
infra-plan: infra-init infra-checkvars
	@pushd terraform; terraform plan -var-file=settings/bastion.tfvars; popd

.PHONY: infra-apply
infra-apply: infra-init infra-checkvars
	@pushd terraform; terraform apply -var-file=settings/bastion.tfvars; popd

.PHONY: infra-destroy
infra-destroy: infra-init infra-checkvars
	@pushd terraform; terraform destroy -var-file=settings/bastion.tfvars; popd

.PHONY: help
help: ## This message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

