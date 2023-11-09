SIMULATOR_IMAGE ?= controlplane/simulator

simulator-dev-image:
	docker build -t $(SIMULATOR_IMAGE):dev -f dev.Dockerfile .

simulator-image: simulator-dev-image
	docker build -t $(SIMULATOR_IMAGE) .

simulator-cli:
	go build -v -o bin/simulator internal/cmd/main.go
