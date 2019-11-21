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
  "Task 1":
    sort-order: 1
    hints:
      - text: "blah blah blah"
        penalty: -5
      - text: |
              compromise the frobnicator
        available-after: 10m
  "Task 2":
    sort-order: 2
    hints:
      - text: "hint 1"
      - text: "hint 2"
```

* Should we merge CHALLENGE.txt into the scenario yaml?
* We could split the challenge text so the tasks each have their own descriptions

### Starting point

Where we dump the user when they do `simulator ssh attack` and type `start`


**First draft in MANIFEST.yaml**
```YAML
attack: container
instance: slave
container-name: nginx
```

**Proposed specification in hints.yaml**

Default behaviour:
```YAML
starting-point:
  mode: attack
```

Land user sshed onto an internal VM instance that is not part of the kubernetes
cluster but is in the private subnet:
```YAML
starting-point:
  mode: internal-instance
  kubectl-access: true
```

Land user sshed onto a node:
```YAML
starting-point:
  mode: node
  node-id: One of <master-{0..n}|node-{0..n}>
  as-root: true
```

Land user exec'ed into a pod:
```YAML
starting-point:
  mode: pod
  pod-name: "compromised-pod"
  pod-namespace: "default"
```
Do not put the user anywhere. This useful for harder scenarios where the whole
cluster may be involved:
```YAML
starting-point:
  mode: null
```

### Kubernetes Kind

Not enough scenarios to specify platforms or versions or requiremnts etc.  Let's revisit when we design more scenarios:

```YAML
local: <True|False>
```
