SIMULATOR_IMAGE ?= controlplane/simulator

lint:
	golangci-lint run -c .golangci.yml

simulator-dev-image: lint
	docker build -t $(SIMULATOR_IMAGE):dev -f dev.Dockerfile .

simulator-image: simulator-dev-image
	docker build -t $(SIMULATOR_IMAGE) .

simulator-cli: lint
	go build -v -o bin/simulator cmd/simulator/main.go
