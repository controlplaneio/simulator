# Standalone simulator

A distributed systems and infrastructure simulator for attacking and debugging Kubernetes

## Usage

The quickest way to get up and running is to simply:

```
make run
```

This will drop you into a bash shell in a launch container.  You will have a program on the `$PATH` named `simulator`
to interact with.  Documentation for using the `simulator` program can be found [in the CLI folder](./cli).

## Roadmap

There is a [roadmap](./docs/roadmap.md) outlining current and planned work.

## Architecture

### [Launching a scenario](./docs/launch.md)

### *TODO* [Validating a scenario](./docs/validation.md)

### *TODO* [Evaluating  scenario progress](./docs/evaluation.md)

### Components

* [Simulator CLI tool](./cli) - Runs in the launch container and orchestrates everything
* [Launch container](./Dockerfile) - Isolates the scripts from the host
* [Terraform Scripts for infrastructure provisioning](./terraform) - AWS infra
* [Perturb.sh](./simulation-scripts/perturb.sh) - sets up a scenario on a cluster
* [Scenarios](./simulation-scripts/scenario)
* *TODO* [Attack container](./Dockerfile.attack) - Runs on the bastion providing all the tools needed to attack a
cluster in the given cloud provider

