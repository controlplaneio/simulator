# Simulator Scenario Development

Configure the Simulator CLI to run in development mode; `simulator config --dev`. When the Simulator CLI launches the
Docker container, it will bind mount the following directories, allowing you to make changes locally and run them,
without having to rebuild the full container image every time.

* ansible
* packer
* terraform

## Scenarios

The scenarios folder contains the [catalogue](../scenarios/scenarios.yaml) listing all the available scenarios. Updating
this with details of the new scenario will add the scenario to the listings available when running
`simulator scenario list` and `simulator scenario describe`.

Create a new directory here, using the scenario id for the name, to add any supplementary files for the scenario e.g.
the learning guide (README.md), the scenario metadata (tasks.yaml), directories for any containerised applications to
support the scenario, and any files to describe the solution.

[//]: # (TODO: determine future of tasks.yaml, format and naming)

## Ansible Configuration

The cluster configuration for the scenario is defined in Ansible using convention over configuration. Each scenario
requires

* An Ansible Playbook to be defined in the [playbooks](../ansible/playbooks) directory, using the scenario id as the
  name of the file.
* An Ansible Role to be define in the [roles](../ansible/roles) directory, using the scenario id as the name of the
  role.

There are also supplementary roles in the [roles](../ansible/roles) directory to support common functionality e.g.
cluster network installation, configuring the player's starting point, configuring socat, etc. These will be added to
over time, as common behaviour is identified.
