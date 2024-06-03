# Creating a Simulator Scenario

Scenarios are a self-contained set of one or more tasks that together comprise a user journey through a Kubernetes issue. Scenarios tell the story of a breached cluster, vulnerability or misconfiguration and guide the user through compromise, escalation of privilege and eventually fixing the cluster. An example of a scenario is the `network-hedgehog-defence` scenario, which explores a compromised database and includes two tasks; access the database across another namespace and add a Kubernetes resource to prevent this.

The objective of this page is to document all the steps needed to create a new scenario for the open-source Kubernetes Simulator.

## Before you start

Determine the type of challenge you want to create and make sure it's not already covered by an existing scenario. A good challenge should have at least a compromise step and a mitigation or fixing step. Ensure that your scenario does not already exist in [the list of scenarios](../simulation-scripts/scenario/). You can read all the challenge texts and hints in each scenario's directory. Familiarise yourself with the structure.

Think about how difficult your scenario will be. Guidelines can be found [here.](./difficulty.md)

## Creating scenarios

### How to add a scenario

1. Create a directory in simulation-scripts/scenario/ for your scenario. This directory should be the name of your scenario and follow the naming convention of the rest of the scenarios of <challenge-topic>-<military-tactic>. You can choose a military tactic at random.

2. Populate your scenario directory as per [the scenario folder structure](#scenario-folder-structure)
In simulation-scripts > scenario create the following structure of folder

3. Put the scenario spec in the [scenario yaml ](../simulation-scripts/scenarios.yaml) and add it to the [README](../README.md). The scenarios in the scenario yaml here will be launchable from the launch container on running the simulator.

### What a good scenario contains

A good scenario should be focused on one or two related vulnerabilities or misconfigurations. Some examples of current scenarios include:

_Resource definitions_
- secrets in environment variables
- missing pod security policies
- missing network definitions
- overly permissible service accounts

_Core Components_
- misconfigured ports on the kubelet
- no authorization on the kubernetes API
- authentication-less etcd

Additionally there should be a structured story to your scenario. Most of the scenarios include an attack task from for example a pod/container or external server and a mitigation task to fix the misconfiguration. Don't forget to provide all necessary information to the user in the challenge.txt file - if the user needs to attack a port on the node, they will need to be aware of the node IP from the start.

### Scenario Folder structure

![Scenario Structure](./scenario-structure.png)

Scenarios live in the `simulator/simulation-scripts/scenario/` directory. Minimally they must contain:

- _scenario name_: Substitute "scenario-name" with the name of your scenario.
- _apply/scenario.yaml_: This file specifies all Kubernetes resources such as namespaces, deployments and services which are part of the scenario. This Yaml will be deployed when the scenario is started. [Example scenario.yaml](https://github.com/kubernetes-simulator/simulator/blob/master/simulation-scripts/scenario/network-hedgehog-defence/apply/scenario.yaml).
- _challenge.txt_: This file contains metadata for the scenario that you want to make available to the user, such as scenario description, starting point, difficulty level and the descriptions of the tasks. This will get displayed on entering the scenario as the challenge text to the user. The _challenge.txt_ can be templated by using any of the environment variables available to the attack container, for example `$MASTER_IP_ADDRESSES`. Also, `##IP` can be prefixed to a Deployment, Pod or DaemonSet name to insert the pod IPs for their associated pods. For example, `##IPfrontend` would be substituted with pod IPs for all pods with names containing 'frontend'. `##NAME` will insert the full pod names and `##HIP` will insert the host IPs of the pods. Example [challenge.txt](https://github.com/kubernetes-simulator/simulator/blob/master/simulation-scripts/scenario/network-hedgehog-defence/challenge.txt).
- _tasks.yaml_: This file contains the task and hint specs for all tasks in the scenario. It also determines where you start - e.g. in a pod on the cluster, on a node or even outside the cluster. [Example tasks.yaml](https://github.com/kubernetes-simulator/simulator/blob/master/simulation-scripts/scenario/network-hedgehog-defence/tasks.yaml). To see the required structure and options of the tasks.yaml, see the guide [here](https://github.com/kubernetes-simulator/simulator/blob/master/docs/tasks-yaml-format.md).

Additionally the `simulator/simulation-scripts/scenario/` directory can include a number of supporting bash scripts which can we executed on the master and worker nodes. These files help with particular setups on the nodes such as installing tooling, creating files or setting up users. The scripts are executed on startup and must have the following names:

- _worker-any.sh_: runs on a randomly chosen worker node
- _worker-1.sh_: runs on worker node 1
- _worker-2.sh_: runs on worker node 2
- _workers-every.sh_: runs on worker nodes 1 & 2
- _nodes-every.sh_: runs on master and worker nodes 1 and 2
- _master.sh_: runs on master

Any other names will cause an error. See this example [master.sh](https://github.com/kubernetes-simulator/simulator/blob/master/simulation-scripts/scenario/etcd-inverted-wedge/master.sh).
