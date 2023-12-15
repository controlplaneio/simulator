# Standalone simulator

`Simulator` is an open source distributed systems and infrastructure simulator for attacking and debugging Kubernetes.

## Download

Please download the latest release from the [releases page](https://github.com/controlplaneio/simulator/releases)

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Terraform](https://www.terraform.io/downloads.html)
- [Packer](https://www.packer.io/downloads)

## Usage

1. Build the container image for the simulator

```bash
  make simulator-image
```

2. Build the AMIs
3. Configure the cli
4. Provision the infrastructure
5. Install the scenario
