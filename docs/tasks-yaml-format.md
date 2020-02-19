# Tasks YAML file Format

## Top Level

There are five top level fields in the spec:

```YAML
category: sample
difficulty: Easy
objective: Sample yaml
kind: cp.simulator/scenario:<semver-version>
tasks: 
  ...
```

* `category`is the broad subject of the scenario, for example one might be RBAC. 
* The `difficulty` describes how hard the scenario is. 
* `objective` is a short sentence about the goals of the scenario. 
* `kind` is the version of tasks.yaml used in this file.
* `tasks:` is explained in more detail below:

## Tasks and Hints

An example of a tasks stanza can be seen below:

```YAML
tasks:
  "1":
    sortOrder: 1
    hints:
      - text: "blah blah blah"
        penalty: 5
      - text: compromise the frobnicator
        penalty: 10
  "2":
    sortOrder: 2
    hints:
      - text: "hint 1"
        penalty: 10
      - text: "hint 2"
        penalty: 15
```

* Each task is identified by its name which **must** be a number enclosed in quotes.
* The `sortOrder` is used to sort the tasks for more advanced functions. For now each task name matches with its `sortOrder`
* The `hints` stanza contains two fields for the hints themselves and the penalty for viewing a hint:
  * Text is a string value displayed to the user when the use `next_hint`
  * Penalty is the value subtracted from the users score when they use `next_hint`


## Starting point

The `startingPoint` is a stanza within each task stanza that controls where the user is placed within the cluster when they start a task and run `starting_point` from within the attack container. `startingPoint` has several possible modes:

Do not put the user anywhere. This useful for harder scenarios where the whole cluster may be involved:

```YAML
startingPoint:
  mode: attack
```

User ssh's onto an internal VM instance that is not part of the kubernetes cluster but is in the private subnet:

```YAML
startingPoint:
  mode: internal-instance
  kubectlAccess: true
```

* `kubectlAccess` can be set to `true` or `false` to control whether the user can use `kubectl` to reach the cluster.

User ssh's onto a node:

```YAML
startingPoint:
  mode: node
  nodeId: One of <master-{0..n}|node-{0..n}>
  asRoot: true
```

User exec's into a specific pod:

```YAML
startingPoint:
  mode: pod
  podName: compromised-pod
  podNamespace: default
```

The pod starting point can also use two optional fields:

* `containerName` to choose a specific container in a pod to start in. This is required for multi-container pods.
* `podHost` to choose a pod on a specific host. Options are one of `master-0`, `node-0` or `node-1`. It is recommended to use this option with a `DaemonSet` as it can be guaranteed that a pod exists on your chosen host.

A starting point using these options is below:

```YAML
startingPoint:
  mode: pod
  podName: compromised-pod
  podNamespace: default
  containerName: container-2
  podHost: node-0
```

Default behaviour:

```YAML
startingPoint:
  mode: null
```

* Will also display a warning about an unconfigured `startingPoint`. To suppress this use the functionally identical `mode: attack`.

## Sample `tasks.yaml`

```YAML
category: sample
difficulty: Easy
kind: cp.simulator/scenario:1.0.0
objective: Sample yaml
tasks:
  "1":
    hints:
    - penalty: 10
      text: sample 1
    - penalty: 10
      text: test 1
    - penalty: 10
    sortOrder: 1
    startingPoint:
      mode: pod
      podName: test
      podNamespace: sample
    summary: Sample summary 1
  "2":
    hints:
    - penalty: 10
      text: sample 2
    - penalty: 10
      text: test 2
    sortOrder: 2
    startingPoint:
      kubectlAccess: true
      mode: internal-instance
    summary: Sample summary 2
```
