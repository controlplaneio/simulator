SIMULATOR_IMAGE ?= controlplane/simulator

.PHONY: help
help: ## Show this help message
	@awk 'BEGIN {FS = ":.*?##"} /^[a-zA-Z_-]+:.*?##/ { sub("\\\\", "", $$1); printf "\033[36m%-30s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

lint: ## Lint the code
	golangci-lint run -c .golangci.yml

simulator-dev-image: lint ## Lint the code and build the development Docker image
	docker build -t $(SIMULATOR_IMAGE):dev -f dev.Dockerfile .

simulator-image: simulator-dev-image ## Build the Docker image for the simulator
	docker build -t $(SIMULATOR_IMAGE) .

simulator-cli: lint ## Build the simulator CLI
	go build -v -o bin/simulator cmd/simulator/main.go

build: simulator-dev-image simulator-image simulator-cli ## Build docker images and the CLI binary
