<!--
This file is evaled by a quickly cobbled together bash script to replace the variables.
- Backticks are imterpreted by bash so use <code> for inline code and <pre> for code blocks.
- If you need to include bsah code snippets you will need to change how the templating works.
-->
# Simulator CLI

## Development and build

Development targets are specified in the [Makefile](./Makefile).

The [Dockerfile](./Dockerfile) in this folder is used to isolate and lock down versions for testing and building and not
for publishing.

Publishing of the simulator as a docker image will be [the launch container Dockerfile](../Dockerfile).
Running  <code>make run</code> (or <code>make build</code>) in the root of this repository will implicitly call
<code>make docker-build</code> in this folder.

The remaining make targets can be run inside the build container or outside as you choose for local development.

<pre>
all: help            
build                 golang build
coverage              runs golang unit tests with coverage and opens a browser with the results
dep                   install dependencies for other targets
doc                   generate markdown documentation for packages
help-no-color:       
test-acceptance       run bats acceptance tests
test                  run all tests
test-unit             run golang unit tests
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
