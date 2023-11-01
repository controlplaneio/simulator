SIMULATOR_IMAGE ?= controlplane/simulator

simulator-dev-image:
	docker build -t $(SIMULATOR_IMAGE):dev -f dev.Dockerfile .

simulator-dev-image-run:
	docker run -it --rm \
		-v ~/.aws:/home/ubuntu/.aws \
		-v $(shell pwd)/ansible:/simulator/ansible \
		-v $(shell pwd)/config:/simulator/config:rw \
		-v $(shell pwd)/packer:/simulator/packer \
		-v $(shell pwd)/terraform:/simulator/terraform \
		--env-file aws-environment \
		--entrypoint bash \
		$(SIMULATOR_IMAGE):dev \

simulator-image: simulator-dev-image
	docker build -t $(SIMULATOR_IMAGE) .

simulator-image-run:
	docker run -it --rm \
		-v ~/.aws:/home/ubuntu/.aws \
		-v $(shell pwd)/config:/simulator/config:rw \
		--env-file aws-environment \
		--entrypoint bash \
		$(SIMULATOR_IMAGE):latest \

simulator-cli:
	go build -v -o bin/simulator internal/cmd/main.go


