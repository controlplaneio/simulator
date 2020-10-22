[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/kubernetes-simulator/simulator/blob/master/LICENSE)
[![Platforms](https://img.shields.io/badge/Platform-Linux|MacOS-blue.svg)](https://github.com/kubernetes-simulator/simulator/blob/master/README.md)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/kubernetes-simulator/simulator/graphs/commit-activity)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![GoDoc](https://godoc.org/github.com/controlplaneio/simulator-standalone?status.svg)](https://godoc.org/github.com/controlplaneio/simulator-standalone)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubernetes-simulator/simulator)](https://goreportcard.com/report/github.com/kubernetes-simulator/simulator)

# Simulator

A distributed systems and infrastructure simulator for attacking and debugging
Kubernetes: <code>simulator</code> creates a Kubernetes cluster for you in your AWS
account; runs scenarios which misconfigure it and/or leave it vulnerable to
compromise and trains you in mitigating against these vulnerabilities.

For details on why we created this project and our ultimate goals take a look at the [vision statement](./docs/vision-statement.md).

## Before you start

## Usage

Ensure the [AWS requirements](#aws-configuration) are met and configured.

Clone this repository and run:

<pre>
make run
</pre>

![Simulator startup](./docs/simulator.png)

This will drop you into a bash shell in a launch container.  You will have a
program on the <code>PATH</code> named <code>simulator</code> to interact with.

## Simulator CLI Usage

### Creating Environment And Lauching Scenario

Before you launch your environment, please see the **How It All Works** section to ensure you have the necessary credentials and permissions in place and know what you are standing up.

_Create a remote state bucket for terraform_
<pre>
simulator init
</pre>
You will be asked for the name for a S3 bucket for the Terraform remote state, which must be globally unique as per AWS standards.  If this does not exist it will be created, otherwise the existing bucket will be used.

_Create the infra if it isn't there_
<pre>
simulator infra create
</pre>
This will standup the infrastructure, including an initial Kubernetes Cluster

_List available scenarios_
<pre>
simulator scenario list
</pre>
This will list all currently available scenarios.  You can filter the list by difficulty or category using the appropriate arguments

_To get a better idea of what is involved in each scenario, e.g. node-shock-tactics_
<pre>
simulator scenario describe node-shock-tactics
</pre>

_Launch a scenario (sets up your cluster)_
<pre>
simulator scenario launch node-shock-tactics
</pre>
This will launch your selected scenario.

_Login to the environment_
<pre>
simulator ssh attack
</pre>

Running <code>simulator ssh attack</code> logs you into a container running on the Bastion host.  Upon login, an outline of the challenge will be displayed.  In addition, shortcuts for logging into the master, or nodes, of the Kubernetes cluster and how to show hints and start tasks are displayed.  <code>starting_point</code> logs you into the correct starting point for the task you have started in the scenario.

From within the ssh attack container you have access to a range of helper commands such as <code>start_task</code>, <code>next_hint</code> and <code>show_hints</code>.

![Bastion container initial login](./docs/bastion.png)

_Start task_
<pre>
start_task 1
</pre>

The <code>start_task</code> command is used to inform the simulator which task you are undertaking with the number associated, and therefore what hints are available to aid you with that task.

_Accessing hints_
<pre>
next_hint
</pre>

The <code>next_hint</code> command will provide a hint to help you complete the task you have started with the <code>start_task</code> command.

_Viewing all hints that have been requested_
<pre>
show_hints
</pre>

The <code>show_hints</code> command will display all the hints you have requested to that point, in the task you have started.

_Ending a task_
<pre>
end_task
</pre>

or

<pre>
start_task 2
</pre>

The <code>end_task</code> command will end the task you are currently on and score you. This will also happen if you move onto the a new task with the <code>start_task</code> command. When you end a task in one of these ways, you will be scored 100 points minus the amount of hints you have displayed for that task at a value of -10 each. This is defined in the scenario task manifest.

### Cleaning Up Environment

_Destroy your cluster when you are done_
<pre>
simulator infra destroy
</pre>
Once you have finished you should destroy the environment to ensure that no additional costs are incurred.  You run this command in the launch container.

### Scenarios

The following scenarios are currently shipped with the simulator:

<pre>
     container-ambush
     container-defeat-in-detail
     container-phalanx-formation
     etcd-inverted-wedge
     master-shell-scrape
     master-encirclement
     network-feint
     network-hammer-and-anvil
     network-hedgehog-defence
     network-swarming
     node-amphibious-operations
     node-raiding
     node-shock-tactics
     policy-echelon-formation
     policy-fire-support
     policy-force-dispersal
     policy-vertical-envelopment
     rbac-contact-drill
     rbac-sangar
     rbac-flanking-maneuver
     rbac-shoot-and-scoot
     secret-high-ground
     secret-tank-desant
</pre>

But you can write your own and we really welcome any contributions of new scenarios :)

## How It All Works

### Infrastructure Design

![Terraform AWS infrastructure](./terraform/docs/aws-bastion-host-1.png)

### AWS Configuration

Simulator uses terraform to provision its infrastructure.  Terraform in turn
honours [all the rules for configuring the AWS
CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html#cli-environment).
Please follow these.

You can provide your credentials via the <code>AWS_ACCESS_KEY_ID</code> and
<code>AWS_SECRET_ACCESS_KEY</code> environment variables, containing your AWS Access Key
and AWS Secret Key respectively. Note that setting your AWS credentials using
either these (or legacy) environment variables will override the use of
<code>AWS_SHARED_CREDENTIALS_FILE</code> and <code>AWS_PROFILE</code>. The <code>AWS_DEFAULT_REGION</code> and
<code>AWS_SESSION_TOKEN</code> environment variables are also used, if applicable.

https://www.terraform.io/docs/backends/types/s3.html

**All the <code>AWS_*</code> configuration environment variables you have set will be propagated into the container**

### Troubleshooting AWS

- If you get a timeout when running <code>simulator infra create</code> after about 10 minutes, the region you are using
is probably running slowly.  You should run <code>simulator infra destroy</code> and then retry <code>simulator infra
create</code>
- <code>AWS_REGION</code> vs <code>AWS_DEFAULT_REGION</code> - There have been
[some issues](https://github.com/aws/aws-sdk-go/issues/2103) with the
[Go AWS client region configuration](https://github.com/aws/aws-sdk-go#configuring-aws-region)
- [Multi-account](https://www.terraform.io/docs/backends/types/s3.html#multi-account-aws-architecture)

### Terraform

Refer to the [simulator terraform documentation](./terraform/README.md)

### SSH

Simulator, whether run in the launch container or on the host machine, will generate its own SSH RSA key pair.  It will configure the cluster to allow access only with this keypair and automates writing SSH config and keyscanning the bastion on your behalf using custom SSH config and known_hosts files.  This keeps all simulator-related SSH config separate from any other SSH configs you may have. All simulator-related SSH files are written to <code>~/.kubesim</code> and are files starting <code>cp_simulator_</code>

**If you delete any of the files then simulator will recreate them and reconfigure the infrastructure as necessary on the
next run**

### Scenario Launching (perturb)

Refer to the [scenario launch documentation](./docs/scenario.md)

### Launch Flow Sequence

Refer to the [Launch scenario flow documentation](./docs/launch.md)

### Scenario Task Definition

Refer to the [Tasks YAML file Format documentation](./docs/tasks-yaml-format.md)

### Create a new scenario

Refer to the [Create Scenario documentation](./docs/create-scenario.md)

### The kubesim script

<code>kubesim</code> is a small script written in BASH for getting users up and running with simulator as fast as possible. It pulls the latest version of the simulator container and sets up some options for running the image. It can be installed with the following steps:

<pre>
# cURL the script from GitHub
curl -Lo kubesim https://github.com/kubernetes-simulator/simulator/releases/latest/download/kubesim
# Make it executeable
chmod a+x kubesim
# Place the script on your path
cp kubesim /usr/local/bin
</pre>

## Development Workflow

### Git hooks

To ensure all Git hooks are in place run the following:

<pre>
make setup-dev
</pre>

Development targets are specified in the [Makefile](./Makefile).

You can see all the available targets with descriptions by running `make help`

### Git commits

We follow [the conventional commit specification](https://www.conventionalcommits.org/en/v1.0.0-beta.4/).
Please ensure your commit messages adhere to this spec.

## Architecture

### [Launching a scenario](./docs/launch.md)

### *TODO* [Validating a scenario](./docs/validation.md)

### *TODO* [Evaluating  scenario progress](./docs/evaluation.md)

### Components

* [Simulator CLI tool](./cmd) - Runs in the launch container and orchestrates everything
* [Launch container](./Dockerfile) - Isolates the scripts from the host
* [Terraform Scripts for infrastructure provisioning](./terraform) - AWS infrastructure
* [Perturb.sh](./simulation-scripts/perturb.sh) - Sets up a scenario on a cluster
* [Scenarios](./simulation-scripts/scenario) - The scenario setup scripts and definitions
* [Attack container](./attack) - Runs on the bastion providing all the tools needed to attack a
cluster in the given cloud provider

### Specifications

* [tasks.yaml format](./docs/tasks-yaml-format.md)

* If you need to make changes to the format you should update this documentation.
* Any changes should be accompanied by a bump of the version in the <code>kind</code>
property
* Use the <code>migrate-hints</code> devtool to update the existing scenarios en-masse.
You can make this tool available on your PATH by running <code>make devtools</code>

### Simulator API Documentation

The simulator API docs are available on godoc.org:

* [Scenario](https://godoc.org/github.com/controlplaneio/simulator-standalone/pkg/scenario)
* [Simulator](https://godoc.org/github.com/controlplaneio/simulator-standalone/pkg/simulator)
* [ChildMinder](https://godoc.org/github.com/controlplaneio/simulator-standalone/pkg/childminder)
* [SSH](https://godoc.org/github.com/controlplaneio/simulator-standalone/pkg/ssh)
* [Util](https://godoc.org/github.com/controlplaneio/simulator-standalone/pkg/util)

## Contributing

Guidelines for contributors, and instructions on working with specific
parts of Simulator are in the [Contributor Guide](./CONTRIBUTING.md).
