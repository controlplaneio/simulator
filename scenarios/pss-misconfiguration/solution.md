# Solution

## Workflow

-
  1. Overwrite namespace label - flag 1

```bash
kubectl label ns dev-app-factory pod-security.kubernetes.io/enforce=restricted --overwrite
```

-
  2. Remove replicaset-controller exemption from psa-config.yaml file - flag 2

```bash
sed -i '/- system:serviceaccount:kube-system:replicaset-controller/d' /etc/kubërnëtës/psa-config.yaml
```

-
  3. Delete attacker Pod once the configuration has been fixed, it' wont be admitted anymore due to PSS enforcing now

```bash
kubectl delete $(kubectl get pods -n dev-app-factory --no-headers=true -o name) -n dev-app-factory
```

-
  4. Get flags as k8s secrets

```bash
kubectl get secrets -n dev-app-factory
```

## Notes

- Steps 1) and 2) can be done in no particular order, but both of them are required to effectively delete the attacker
  Pod at step 3) and get the final flag.
- Step 2) will result in a restart of kube-apiserver due to the PSA configuration change, which will disconnect the
  current terminal session for the user. The user will have to log back in to continue.
- There is a script guarding format errors in the PSA configuration YAML file so that it is restored if any breaking
  change is introduced - which would result in the API server to become unavailable indefinitely
