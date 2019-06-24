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

### Infrastructure

<pre>
${infra_help}
</pre>

## API Documentation

* [Scenario](./docs/scenario.md)
