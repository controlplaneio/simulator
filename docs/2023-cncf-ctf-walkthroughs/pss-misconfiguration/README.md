# Scenario: PSS Misconfiguration

- [Learning Outcomes](#learning-outcomes)
- [Challenge Description](#challenge-description)
- [Guided Walkthrough](#guided-walkthrough)
  - [Step 1: Enumerate resources and discover the PSS labels](#step-1-enumerate-resources-and-discover-the-pss-labels)
  - [Step 2: Patching the PSS labels](#step-2-patching-the-pss-labels)
  - [Step 3: Discovery of the Pod Security Admission resource and exemptions](#step-3-discovery-of-the-pod-security-admission-resource-and-exemptions)
  - [Step 4: Acquiring the remaining flags and remove the adversary workload](#step-4-acquiring-the-remaining-flags-and-remove-the-adversary-workload)
- [Remediation and Security Considerations](#remediation-and-security-considerations)

## Learning Outcomes

The purpose of PSS Misconfiguration is to provide participants a deeper understanding of how Pod Security Standards and Pod Security Admission work. This includes:

1. Pod Security levels and how they are configured
2. The different Pod Security Admission label options
3. How Pod Security Admission exemptions function

## Challenge Description

```
The platform team spent their Christmas holidays migrating from pod security policies to the new fancy standard, locking down the k8s fleet!

Nonetheless, an attacker stole a developer credential and has managed to run his own unconstrained Pod in the dev-app-factory namespace, to later pivot on the host.

How is that even possible?! Can you find a way to remediate this and ensure that the attacker Pod can't run anymore in the cluster?
```

## Guided Walkthrough

### Step 1: Enumerate resources and discover the PSS labels

The first step is to enumerate resources we have access to. Let's start by finding out what our service account has permissions to perform on the cluster.

```bash
root@admin:~# kubectl auth can-i --list
Resources                                       Non-Resource URLs                     Resource Names      Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []                  [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []                  [create]
limitranges                                     []                                    [dev-app-factory]   [get list patch update]
namespaces                                      []                                    [dev-app-factory]   [get list patch update]
resourcequotas                                  []                                    [dev-app-factory]   [get list patch update]
                                                [/.well-known/openid-configuration]   []                  [get]
                                                [/api/*]                              []                  [get]
                                                [/api]                                []                  [get]
                                                [/apis/*]                             []                  [get]
                                                [/apis]                               []                  [get]
                                                [/healthz]                            []                  [get]
                                                [/healthz]                            []                  [get]
                                                [/livez]                              []                  [get]
                                                [/livez]                              []                  [get]
                                                [/openapi/*]                          []                  [get]
                                                [/openapi]                            []                  [get]
                                                [/openid/v1/jwks]                     []                  [get]
                                                [/readyz]                             []                  [get]
                                                [/readyz]                             []                  [get]
                                                [/version/]                           []                  [get]
                                                [/version/]                           []                  [get]
                                                [/version]                            []                  [get]
                                                [/version]                            []                  [get]
```

It looks like we have access to `namespaces`, `limitranges` and `resourcequotas` in the `dev-app-factory` namespace.

Let's check what else we have access to in that namespace.

```bash
root@admin:~# kubectl auth can-i --list -n dev-app-factory
Resources                                       Non-Resource URLs                     Resource Names      Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []                  [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []                  [create]
limitranges                                     []                                    [dev-app-factory]   [get list patch update]
namespaces                                      []                                    [dev-app-factory]   [get list patch update]
resourcequotas                                  []                                    [dev-app-factory]   [get list patch update]
pods                                            []                                    []                  [get list watch delete]
secrets                                         []                                    []                  [get list watch]
daemonsets.apps                                 []                                    []                  [get list watch]
deployments.apps                                []                                    []                  [get list watch]
statefulsets.apps                               []                                    []                  [get list watch]
                                                [/.well-known/openid-configuration]   []                  [get]
                                                [/api/*]                              []                  [get]
                                                [/api]                                []                  [get]
                                                [/apis/*]                             []                  [get]
                                                [/apis]                               []                  [get]
                                                [/healthz]                            []                  [get]
                                                [/healthz]                            []                  [get]
                                                [/livez]                              []                  [get]
                                                [/livez]                              []                  [get]
                                                [/openapi/*]                          []                  [get]
                                                [/openapi]                            []                  [get]
                                                [/openid/v1/jwks]                     []                  [get]
                                                [/readyz]                             []                  [get]
                                                [/readyz]                             []                  [get]
                                                [/version/]                           []                  [get]
                                                [/version/]                           []                  [get]
                                                [/version]                            []                  [get]
                                                [/version]                            []                  [get]
```

We can see that we have access to `pods`, `secrets`, `daemonsets.apps`, `deployments.apps` and `statefulsets.apps` in the `dev-app-factory` namespace. Let's see what workloads and secrets are currently deployed.

```bash
root@admin:~# kubectl get pods -n dev-app-factory
NAME                     READY   STATUS    RESTARTS   AGE
nginx-759f854c9c-zcrxb   1/1     Running   0          6m48s

root@admin:~# kubectl get deployments -n dev-app-factory
NAME    READY   UP-TO-DATE   AVAILABLE   AGE
nginx   1/1     1            1           12m

root@admin:~# kubectl get daemonsets -n dev-app-factory
No resources found in dev-app-factory namespace.

root@admin:~# kubectl get statefulsets -n dev-app-factory
No resources found in dev-app-factory namespace.

root@admin:~# kubectl get secrets -n dev-app-factory
No resources found in dev-app-factory namespace.
```

There is pod `nginx-759f854c9c-zcrxb` running which must be the unconstrained workload that the adversary has deployed. It seems to be configured via a Kubernetes deployment. We have permissions to `delete` pods, so let's remove it.

```bash
root@admin:~# kubectl delete pods nginx-759f854c9c-zcrxb -n dev-app-factory
pod "nginx-759f854c9c-zcrxb" deleted
```

Well that was easy...or was it?

```bash
root@admin:~# kubectl get pods -n dev-app-factory
NAME                     READY   STATUS    RESTARTS   AGE
nginx-759f854c9c-pjgxl  1/1     Running   0          13s
```

It seems that the adversary is using a deployment for the workload and we cannot remove that. There must be another way to apply an enforcement policy to the namespace? This is where Pod Security Standards come in. Let's inspect the namespace further by using the `-o` flag to see whether any Pod Security Standards have been applied.

```bash
root@admin:~# kubectl get namespaces dev-app-factory -oyaml
```

```yaml
apiVersion: v1
kind: Namespace
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","kind":"Namespace","metadata":{"annotations":{},"labels":{"pod-security.kubernetes.io/enforce":"privileged","pod-security.kubernetes.io/warn":"restricted"},"name":"dev-app-factory"}}
  creationTimestamp: "2023-11-22T14:48:49Z"
  labels:
    kubernetes.io/metadata.name: dev-app-factory
    pod-security.kubernetes.io/enforce: privileged
    pod-security.kubernetes.io/warn: restricted
  name: dev-app-factory
  resourceVersion: "771"
  uid: 7867345d-f445-44cb-9c1b-ef6e17bfa815
spec:
  finalizers:
  - kubernetes
status:
  phase: Active
```

We can see that the namespace has two labels, `pod-security.kubernetes.io/enforce: privileged` and `pod-security.kubernetes.io/warn: restricted`. These labels are used by the Pod Security Standards to determine what level of security a pod should be enforced and the warning message returned if the pod does not meet the requirements.

> Note: It is highly recommended to review the [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/) documentation to understand how the labels work and what controls they apply.

### Step 2: Patching the PSS labels

The current configuration enforces `privileged` which provides an **"unrestricted"** policy allowing pods to run with any privilege but the warning message is set to `restricted` which potentially is misleading to any one deploying workloads into the namespace. Let's correct this.

Based on the [PSS Configuration Documentation](https://kubernetes.io/docs/tasks/configure-pod-container/enforce-standards-namespace-labels/#add-labels-to-existing-namespaces-with-kubectl-label) we can use the `kubectl label` command to patch the namespace labels.

```bash
root@admin:~# kubectl label namespace dev-app-factory pod-security.kubernetes.io/enforce=restricted --overwrite
Warning: existing pods in namespace "dev-app-factory" violate the new PodSecurity enforce level "restricted:latest"
Warning: nginx-759f854c9c-pjgxl: allowPrivilegeEscalation != false, unrestricted capabilities, runAsNonRoot != true, runAsUser=0, seccompProfile
namespace/dev-app-factory labeled
```

Great that has worked. Interesting, we have a warning about the nginx workload in the `dev-app-factory` namespace as it is running without a number of security contexts and a seccomp profile. Now we've completed an remediation action, let's see if anything has changed?

```bash
root@admin:~# kubectl get secrets -n dev-app-factory
NAME       TYPE     DATA   AGE
flag-xyz   Opaque   1      18m
```

It looks like we got our first flag, let's see what it is.

```bash
root@admin:~# kubectl get secrets flag-xyz -n dev-app-factory -ojson
{
    "apiVersion": "v1",
    "data": {
        "flag": "ZmxhZ19jdGZ7TUlTQ09ORklHX0dPVF9QU0FfTEVWRUxfV1JPTkd9"
    },
    "kind": "Secret",
    "metadata": {
        "creationTimestamp": "2023-11-22T15:16:48Z",
        "name": "flag-xyz",
        "namespace": "dev-app-factory",
        "resourceVersion": "3199",
        "uid": "57c9a2ba-7505-45e6-b04e-436241e99e6c"
    },
    "type": "Opaque"
}

root@admin:~# echo "ZmxhZ19jdGZ7TUlTQ09ORklHX0dPVF9QU0FfTEVWRUxfV1JPTkd9" | base64 -d
flag_ctf{MISCONFIG_GOT_PSA_LEVEL_WRONG}
```

Now with the PSS restriction in place, let's try and remove the workload.

```bash
root@admin:~# kubectl delete pods nginx-759f854c9c-pjgxl -n dev-app-factory
pod "nginx-759f854c9c-pjgxl" deleted
```

We've managed to delete the pod but let's check again...

```bash
root@admin:~# kubectl get pods -n dev-app-factory
NAME                     READY   STATUS    RESTARTS   AGE
nginx-759f854c9c-xr9dl   1/1     Running   0          29s
```

Nope, it seems to have been redeployed again, even with the PSS policy in place. Potentially there is another level of persistence going on here, let's see what else we can find.

### Step 3: Discovery of the Pod Security Admission resource and exemptions

If we enumerate the `admin` file system we cannot find any other service account tokens.

```bash
root@admin:~# ls -la /var/run/secrets/kubernetes.io/serviceaccount/
total 4
drwxrwxrwt 3 root root  140 Nov 22 15:33 .
drwxr-xr-x 3 root root 4096 Nov 22 15:33 ..
drwxr-xr-x 2 root root  100 Nov 22 15:33 ..2023_11_22_15_33_12.141272112
lrwxrwxrwx 1 root root   31 Nov 22 15:33 ..data -> ..2023_11_22_15_33_12.141272112
lrwxrwxrwx 1 root root   13 Nov 22 15:33 ca.crt -> ..data/ca.crt
lrwxrwxrwx 1 root root   16 Nov 22 15:33 namespace -> ..data/namespace
lrwxrwxrwx 1 root root   12 Nov 22 15:33 token -> ..data/token
```

However, if we look in `/etc` we notice an interesting directory.

```bash
root@admin:~# ls -la /etc
total 276
drwxr-xr-x 1 root root    4096 Nov 22 15:33  .
drwxr-xr-x 1 root root    4096 Nov 22 15:33  ..
-rw------- 1 root root       0 Aug 14 00:00  .pwd.lock
drwxr-xr-x 2 root root    4096 Aug 14 00:00  alternatives
...
drwxr-xr-x 3 root root    4096 Nov 22 15:33  kubernetes
drwxr-xr-x 2 root root    4096 Nov 22 15:33 'kub'$'\303\253''rn'$'\303\253''t'$'\303\253''s'
...
```

This is not a normal directory name, let's see what is inside.

```bash
root@admin:~# ls -la /etc/'kub'$'\303\253''rn'$'\303\253''t'$'\303\253''s'
total 12
drwxr-xr-x 2 root root 4096 Nov 22 15:33 .
drwxr-xr-x 1 root root 4096 Nov 22 15:33 ..
-rw-r--r-- 1 root root  559 Nov 22 15:33 psa-config.yaml
```

We've appeared to have a Pod Security Admission configuration file. Let's see what it contains.

```bash
root@admin:~# cat /etc/'kub'$'\303\253''rn'$'\303\253''t'$'\303\253''s'/psa-config.yaml
```

```yaml
apiVersion: apiserver.config.k8s.io/v1
kind: AdmissionConfiguration
plugins:
- configuration:
    apiVersion: pod-security.admission.config.k8s.io/v1beta1
    defaults:
      audit: restricted
      audit-version: latest
      enforce: restricted
      enforce-version: latest
      warn: restricted
      warn-version: latest
    exemptions:
      namespaces:
      - kube-system
      - platform
      runtimeClasses: []
      usernames:
      - system:serviceaccount:kube-system:replicaset-controller
    kind: PodSecurityConfiguration
  name: PodSecurity
```

Let's breakdown the Pod Security Admission configuration file.

1. We can see the defaults for pod security levels are set to `restricted` with the pod security admission labels of `audit`, `enforce` and `warn`.
2. There are two namespace that are exempt from the pod security admission labels, `kube-system` and `platform`.
3. There is a service account `system:serviceaccount:kube-system:replicaset-controller` that is also exempt from the pod security admission labels.

Let's research what the `replicaset-controller` service account is used for. Based on the [Replication Controller Documentation](https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/) we can see that **"A ReplicationController ensures that a specified number of pod replicas are running at any one time."**. If the service account associated with replicaset controller is exempt from the pod security admission labels, then it is possible that the adversary is using this to redeploy the nginx workload that bypasses the admission policy.

Let's remove the exemption from the `psa-config.yaml` file and see if that resolves the issue.

> Note: You will need to install your preferred text editor on the `jumpbox-terminal` to create the file. This can be achieved with `apt update && apt install -y <text-editor>`

```bash
root@admin:~# vim /etc/'kub'$'\303\253''rn'$'\303\253''t'$'\303\253''s'/psa-config.yaml
root@admin:~# Connection to 172.31.2.253 closed.
Connection to 13.40.53.184 closed.
```

We were disconnected! Let's ssh back and see what has happened.

### Step 4: Acquiring the remaining flags and remove the adversary workload

Logging back into the `admin` terminal we can check and see if our remediation action has worked. Previously we were rewarded with secret, so let's check again.

```bash
root@admin:~# kubectl get secrets -n dev-app-factory
NAME        TYPE     DATA   AGE
flag-xyz    Opaque   1      10m23s
flag-wasd   Opaque   1      5m44s
```

Great, we've got our second flag.

```bash
root@admin:~# kubectl get secrets flag-wasd -n dev-app-factory -ojson
{
    "apiVersion": "v1",
    "data": {
        "flag": "ZmxhZ19jdGZ7TUlTQ09ORklHX0dPVF9QU0FfRVhFTVBUSU9OU19XUk9OR30="
    },
    "kind": "Secret",
    "metadata": {
        "creationTimestamp": "2023-11-22T16:23:54Z",
        "name": "flag-wasd",
        "namespace": "dev-app-factory",
        "resourceVersion": "5133",
        "uid": "03e0afde-0121-4d9d-a5f6-344adeab9d75"
    },
    "type": "Opaque"
}

root@admin:~# echo "ZmxhZ19jdGZ7TUlTQ09ORklHX0dPVF9QU0FfRVhFTVBUSU9OU19XUk9OR30=" | base64 -d
flag_ctf{MISCONFIG_GOT_PSA_EXEMPTIONS_WRONG}
```

Finally let's remove that pesky nginx workload.

```bash
root@admin:~# kubectl delete pods nginx-759f854c9c-tnbgc -n dev-app-factory
pod "nginx-759f854c9c-xr9dl" deleted
```

And double check this pod has not been redeployed.

```bash
root@admin:~# kubectl get pods -n dev-app-factory
No resources found in dev-app-factory namespace.
```

Excellent, we've finally removed the adversary workload. Let's check if there are any more flags to collect.

```bash
root@admin:~# kubectl get secrets -n dev-app-factory
NAME         TYPE     DATA   AGE
flag-gotit   Opaque   1      111s
flag-wasd    Opaque   1      10m
flag-xyz     Opaque   1      5m12s
```

That is our final flag.

```bash
root@admin:~# kubectl get secrets flag-gotit -n dev-app-factory -ojson
{
    "apiVersion": "v1",
    "data": {
        "flag": "ZmxhZ19jdGZ7TUlTQ09ORklHX1BTU19IQVNISkFDS19TQURfQkxPQ0tFRH0="
    },
    "kind": "Secret",
    "metadata": {
        "creationTimestamp": "2023-11-22T16:32:55Z",
        "name": "flag-gotit",
        "namespace": "dev-app-factory",
        "resourceVersion": "5993",
        "uid": "c81eab8a-5c62-4bbb-bed5-ea17363b79e0"
    },
    "type": "Opaque"
}

root@admin:~# echo "ZmxhZ19jdGZ7TUlTQ09ORklHX1BTU19IQVNISkFDS19TQURfQkxPQ0tFRH0=" | base64 -d
flag_ctf{MISCONFIG_PSS_HASHJACK_SAD_BLOCKED}
```

Congratulations, you've completed the PSS Misconfiguration.

## Remediation and Security Considerations

Completing the CTF scenario functions as remediation plan, patching Kubernetes resources and evicting an adversary from the cluster. But as reminder of the important actions performed.

- Ensure that the Pod Security Standard labels are configured correctly and that the `enforce` label is set to `restricted` for all namespaces.
- Ensure that the Pod Security Admission exemptions do not include any service accounts which may bypass the Pod Security Standards admission controls.