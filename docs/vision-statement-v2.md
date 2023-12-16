# Simulator V2 Vision

Simulator V2 is major refresh and architectural restructuring of the Simulator project aimed at simplifying its use for
players, core maintainers and scenario developers. The key drivers were simply;

- Refresh and simplify the codebase
- Democratise scenario development

## Refresh and simplify the codebase

The previous version of Simulator consisted of

- a golang cli
- a nodejs cli
- a bash configuration framework
- tcl tests
- terraform configuration
- two distinct container images

This has been consolidated down to

- a golang cli
- packer configuration
- terraform configuration
- ansible playbooks and roles
- two related container images (player, scenario developer)

The core of Simulator is now a golang CLI, that along with the prerequisites of Docker and an AWS account, is all that
is needed to play the scenarios.

We use
- [Packer](https://www.hashicorp.com/products/packer) to create AMIs with as much static configuration as possible
to reduce startup time.
- [Terraform](https://www.hashicorp.com/products/terraform) to provision the Simulator infrastructure.
- [Ansible](https://www.ansible.com/) to configure the infrastructure for the scenarios.

We bundled these tools and their configuration into a
[container image](https://hub.docker.com/repository/docker/controlplane/simulator/general) so that a player does not
have to install these and manage version compatibility.

## Democratise scenario development

At it's core, Simulator is provisioning and configuring infrastructure. Terraform and Ansible are widely used for these
tasks by people managing infrastructure. Moving the scenario provisioning to be Ansible based simplifies the mental
model, provides the opportunity to build up a library of reusable components for scenario development, and makes it
simpler to implement uninstall functionality for scenarios (coming soon), which means players will not have to destroy
and recreate the infrastructure as they play different scenarios.

Simulator V2 is also a scenario development kit. When run in dev mode, you can create a scenario one Ansible task at
a time, iteratively building it up and verifying every step of the way.
