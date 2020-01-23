# Tasks YAML file Format

### Boilerplate

* YAML version
```YAML
kind: cp.simulator/scenario:<semver-version>
```

* Scenario tags (categories)
```YAML
tags:
 - AWS
 - Debug
```

### Tasks and Hints

e.g.
```YAML
tasks:
  "1":
    sortOrder: 1
    hints:
      - text: "blah blah blah"
        penalty: 5
      - text: |
              compromise the frobnicator
  "2":
    sortOrder: 2
    hints:
      - text: "hint 1"
      - text: "hint 2"
```

* Should we merge CHALLENGE.txt into the scenario yaml?
* We could split the challenge text so the tasks each have their own descriptions

### Starting point

Where we dump the user when they do `simulator ssh attack` and type `start`

**Proposed specification in hints.yaml**

Default behaviour:
```YAML
startingPoint:
  mode: attack
```

Land user sshed onto an internal VM instance that is not part of the kubernetes
cluster but is in the private subnet:
```YAML
startingPoint:
  mode: internal-instance
  kubectlAccess: true
```

Land user sshed onto a node:
```YAML
startingPoint:
  mode: node
  nodeId: One of <master-{0..n}|node-{0..n}>
  asRoot: true
```

Land user exec'ed into a pod:
```YAML
startingPoint:
  mode: pod
  podName: "compromised-pod"
  podNamespace: "default"
```
Do not put the user anywhere. This useful for harder scenarios where the whole
cluster may be involved:
```YAML
startingPoint:
  mode: null
```

### Kubernetes Kind

Not enough scenarios to specify platforms or versions or requirements etc.  Let's revisit when we design more scenarios:

```YAML
local: <True|False>
```
