<!--
This file is evaled by a quickly cobbled together bash script to replace the variables.
- Backticks are imterpreted by bash so use <code> for inline code and <pre> for code blocks.
- If you need to include bsah code snippets you will need to change how the templating works.
-->
# Simulator CLI

## Development and build

Development targets are specified in the [Makefile](./Makefile).

This application is built and tested in a multi-stage build in the parent directory's [launch container dockerfile](../Dockerfile)

The remaining make targets can be run locally.

<pre>
all: help
build                 golang build
dep                   install dependencies for other targets
doc                   generate markdown documentation for packages
docker-build          builds a docker image
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
  infra       Interact with AWS to create, query and destroy the required infrastructure for scenarios
  scenario    Interact with scenarios
  version     Prints simulator version

Flags:
  -h, --help      help for simulator
  -v, --verbose   verbose output

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

Global Flags:
  -v, --verbose   verbose output

Use "simulator scenario [command] --help" for more information about a command.
</pre>

### Infrastructure

<pre>
Interact with AWS to create, query and destroy the required infrastructure for scenarios

Usage:
  simulator infra [command]

Available Commands:
  create      Runs terraform to create the required infrastructure for scenarios
  destroy     Tears down the infrastructure created for scenarios
  status      Gets the status of the infrastructure

Flags:
  -h, --help   help for infra

Global Flags:
  -v, --verbose   verbose output

Use "simulator infra [command] --help" for more information about a command.
</pre>

## API Documentation

* [Scenario](./docs/scenario.md)
* [Runner](./docs/runner.md)
