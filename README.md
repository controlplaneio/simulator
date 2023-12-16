[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/controlplaneio/simulator/blob/master/LICENSE)
[![Platforms](https://img.shields.io/badge/Platform-Linux|MacOS-blue.svg)](https://github.com/controlplaneio/simulator/blob/master/README.md)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-blue.svg)](https://conventionalcommits.org)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/controlplaneio/simulator/graphs/commit-activity)

# Simulator

A distributed systems and infrastructure simulator for attacking and debugging Kubernetes: <code>simulator</code>
creates a Kubernetes cluster for you in your AWS account; runs scenarios which misconfigure it and/or leave it
vulnerable to compromise and trains you in mitigating against these vulnerabilities.

For details on why we created this project and our initial goals take a look at the [vision statement](./docs/vision-statement.md).
For details of the vision and drivers for Simulator V2, take a look at [version two](docs/vision-statement-v2.md)

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [AWS Account](https://aws.amazon.com/free)

## Download

Please download the latest release from the [releases page](https://github.com/controlplaneio/simulator/releases)

## AWS Credentials

Simulator supports the following methods of authentication to provision the AWS infrastructure. Refer to
[AWS IAM Permissions](docs/aws-iam-permissions.md) for details of the required permissions.

- [Environment Variables](https://docs.aws.amazon.com/sdkref/latest/guide/environment-variables.html)
- [Shared Credentials file](https://docs.aws.amazon.com/sdkref/latest/guide/file-format.html)

## Getting Started

- Read the [Player Guide](docs/player-guide.md) to learn how to launch environments, deploy, and play scenarios
- Refer to the [Walkthough Guides](docs/2023-cncf-ctf-walkthroughs) for advice

## Core Components

- [Simulator CLI](docs/cli.md)
- [Simulator Container Images](docs/container-images.md)
- [Simulator AMIs](docs/amis.md)
- [Simulator Infrastructure](docs/infrastructure.md)
- [Simulator Scenarios](docs/scenarios.md)
