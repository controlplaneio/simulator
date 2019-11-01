Hints YAML file Format

### Boilerplate

* YAML version
```
kind: cp.simulator/scenario:<semver-version>
```

* Scenario tags (categories)
```
tags:
 - AWS
 - Debug
```

### Hints

* Let's not put challenge text in yaml because pertub will have to parse it with yq/jq

https://gist.github.com/jondkent/2412170f5feba8a3cfc977fea10e1c26

I think we need the hints to be a "sequence of mappings" (YAML terminology) rather than a "sequence of strings"
e.g.
```
hints:
  - text: "blah blah blah"
    penalty: -5
  - text: |
          compromise the frobnicator
    available-after: 10m
```

### Starting point

Where we dump the user when they do `simulator ssh attack` and type `start`


**First draft in MANIFEST.yaml**
```
attack: container
instance: slave
container-name: nginx
```

**Proposed specification in hints.yaml**

Default behaviour:
```
starting-point: 
  mode: attack
```

Land user sshed onto an internal VM instance that is not part of the kubernetes
cluster but is in the private subnet:
```
starting-point:
  mode: internal-instance
```

Land user sshed onto a node:
```
starting-point:
  mode: node
  node-id: One of <master-{0..n}|node-{0..n}>
```

Land user exec'ed into a pod:
```
starting-point:
  mode: pod
  node-type: One of <master|node>
  pod-name: "compromised-pod"
```

### Kubernetes Kind

Not enough scenarios to specify platforms or versions or requiremnts etc.  Let's revisit when we design more scenarios:

```
local: <True|False>
```
