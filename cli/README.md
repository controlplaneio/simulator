# Simulator CLI

## Development

Development is done via make

<pre>
all: help            
build                 golang build
dep                   install dependencies for other targets
doc                   generate markdown documentation for packages
docker-build          builds a docker image
docker-push           pushes the last build docker image
docker-run            runs the last build docker image
help-no-color:       
test-acceptance       acceptance tests
test-go-fmt           golang fmt check
test                  unit and local acceptance tests
test-unit             golang unit tests
test-unit-verbose     golang unit tests (verbose)
</pre>

## Usage

<pre>

A distributed systems and infrastructure simulator for attacking and
debugging Kubernetes

Usage:
  simulator [command]

Available Commands:
  help        Help about any command
  scenario    Interact with scenarios
  version     Prints simulator version

Flags:
  -h, --help   help for simulator

Use "simulator [command] --help" for more information about a command.
</pre>

### Scenarios

<pre>
Interact with scenarios

Usage:
  simulator scenario [command]

Available Commands:
  launch      Launches a scenario
  list        Lists available scenarios

Flags:
  -h, --help   help for scenario

Use "simulator scenario [command] --help" for more information about a command.
</pre>

## API Documentation

* [Scenario](./docs/scenario.md)
