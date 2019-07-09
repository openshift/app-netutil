# Network utility for application running on Kubernetes

## Table of Contents

- [Network Utility](#network-utility)
- [Quick Start](#quick-start)
	- [Build](#build)
	- [Create Test Pod](#create-test-pod) 

## Network Utility
Network Utility is a program that runs in the sidecar container of each Kubernetes pod, it provides gRPC API for applications running in the same pod (but different container) to query network information associated with pod. Network Utility is written in golang and can be built into container image.


## Quick Start
This section explains an example deployment using netutil in sidecar container.

### Build

1. Compile executable:
```
$ make
```

This builds two binaries under `bin/` dir, one is named `server` that will run in sidecar container, the other is client which is an example of application that queries API server via gRPC call. Currently there is only one gRPC method implemented called `GetCPUInfo`.

2. Build container image:
```
$ make image
```

This will build an image containing `server` and `client` binaries, the image will be used when creating test pod.

3. Clean up:
```
$ make clean
```

This cleans up built binaries and softlinks.

### Create Test Pod

 1. Create pod with sidecar container:
```
$ kubectl create -f deployments/pod.yaml
```
 2. Check for logs of sidecar container:
```
$ kubectl logs testpod -c container-1

I0619 10:05:11.187825      16 main.go:36] starting netutil server at: /etc/netutil/net.sock
I0619 10:05:11.190333      16 main.go:57] netutil server start serving
I0619 10:05:11.915010      16 main.go:77] getting cpuset from path: /proc/11/root/sys/fs/cgroup/cpuset/cpuset.cpus
...
```

5. Check for logs of application container
````
$ kubectl logs testpod -c container-2

I0619 10:05:10.910894      11 client.go:22] starting netutil client
I0619 10:05:11.912938      11 client.go:34] os.Getpid() returns: 11
I0619 10:05:11.916298      11 client.go:39] GetCPUInfo Response from NetUtilServer: 0-35
````
