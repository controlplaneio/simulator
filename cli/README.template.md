# Simulator CLI

## Development and build

Development targets are specified in the [Makefile](./Makefile).

The [./Dockerfile](Dockerfile) in this folder is used to isolate and lock down versions for testing and building and not for publishing.

Publishing of the simulator as a docker image will be [../Dockerfile](the launch container Dockerfile in the root). Running  `make run` (or `make build`) in the root of this repository will implicitly call `make docker-build` in this folder.

The remaining make targets can be run inside the build container or outside as you choose for local development.

<pre>
${make}
</pre>

## Usage

<pre>
${help}
</pre>

### Scenarios

<pre>
${scenario_help}
</pre>

## API Documentation

* [Scenario](./docs/scenario.md)
