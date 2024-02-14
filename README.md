[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/controlplaneio/simulator/blob/master/LICENSE)
[![Platforms](https://img.shields.io/badge/Platform-Linux|MacOS-blue.svg)](https://github.com/controlplaneio/simulator/blob/master/README.md)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-blue.svg)](https://conventionalcommits.org)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/controlplaneio/simulator/graphs/commit-activity)

# Simulator

_🔊 Join [the hosted Simulator waitlist](https://kubesim.io) for private scenarios and training content 🚨_

A distributed systems and infrastructure simulator for attacking and debugging Kubernetes: <code>simulator</code>
creates a Kubernetes cluster for you in your AWS account; runs scenarios which misconfigure it and/or leave it
vulnerable to compromise and trains you in mitigating against these vulnerabilities.

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

## More Info

For details on why we created this project, take a look at the [vision statement](./docs/vision-statement.md) and [V2 redesign](/docs/vision-statement-v2.md).

---

Built with ❤ by [https://control-plane.io/](https://control-plane.io/)
