<!--
This file is evaled by a quickly cobbled together bash script to replace the variables.
- Backticks are imterpreted by bash so use <code> for inline code and <pre> for code blocks.
- If you need to include bsah code snippets you will need to change how the templating works.
-->
# Simulator CLI

## Development and build

Development targets are specified in the [Makefile](./Makefile).

<pre>
run                   Runs the simulator - the build stage of the container runs all the cli tests
docker-build          Builds the launch container
docker-test           Run the tests
infra-init            Initialisation needed before interacting with the infra
infra-checkvars       Check the tfvars file exists before interacting with the infra
infra-plan            Show what changes will be applied to the infrastructure
infra-apply           Apply any changes needed to the infrastructure before running a scenario
infra-destroy         Teardown any infrastructure
dep                   Install dependencies for other targets
build                 Run golang build for the CLI program
test                  run all tests except goss tests
test-acceptance       Run bats acceptance tests for the CLI program
test-unit             Run golang unit tests for the CLI program
coverage              Run golang unit tests with coverage and opens a browser with the results
doc                   Generate documentation
</pre>

## Usage

<pre>

A distributed systems and infrastructure simulator for attacking and
debugging Kubernetes

Usage:
  simulator [command]

Available Commands:
  config      Interact with simulator config
  help        Help about any command
  infra       Interact with AWS to create, query and destroy the required infrastructure for scenarios
  scenario    Interact with scenarios
  version     Prints simulator version

Flags:
  -c, --config-file string   the directory where simulator.yaml can be found
  -h, --help                 help for simulator
  -l, --loglevel string      the level of detail in output logging (default "info")

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
  -c, --config-file string   the directory where simulator.yaml can be found
  -l, --loglevel string      the level of detail in output logging (default "info")

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
  -c, --config-file string   the directory where simulator.yaml can be found
  -l, --loglevel string      the level of detail in output logging (default "info")

Use "simulator infra [command] --help" for more information about a command.
</pre>

### Config

<pre>
Interact with simulator config

Usage:
  simulator config [command]

Available Commands:
  get         Gets the value of a setting

Flags:
  -h, --help   help for config

Global Flags:
  -c, --config-file string   the directory where simulator.yaml can be found
  -l, --loglevel string      the level of detail in output logging (default "info")

Use "simulator config [command] --help" for more information about a command.
</pre>

## API Documentation

* [Scenario](./docs/api/scenario.md)
* [Runner](./docs/api/runner.md)
