<!--
This file is evaled by a quickly cobbled together bash script to replace the variables.
- Backticks are imterpreted by bash so use <code> for inline code and <pre> for code blocks.
- If you need to include bsah code snippets you will need to change how the templating works.
-->
# Simulator CLI

## Development and build

Development targets are specified in the [Makefile](./Makefile).

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

### Config

<pre>
${config_help}
</pre>

## API Documentation

* [Scenario](./docs/api/scenario.md)
* [Runner](./docs/api/runner.md)
