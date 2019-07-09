# Network utility for application running on Kubernetes

## Table of Contents

- [Network Utility](#network-utility)
- [Quick Start](#quick-start)
	- [Build](#build)
	- [Create Test Pod](#create-test-pod) 

## Network Utility
Network Utility is a program that runs in the sidecar container of each Kubernetes pod, it provides gRPC API for applications running in the same pod (but different container) to query network information associated with pod. Network Utility is written in golang and can be built into container image.

Currently there are three gRPC server methods implemented in `NetUtilServer` interface:

```
// NetUtilServer is the server API for NetUtil service.
type NetUtilServer interface {
        GetCPUInfo(context.Context, *CPURequest) (*CPUResponse, error)
        GetEnv(context.Context, *EnvRequest) (*EnvResponse, error)
        GetNetworkStatus(context.Context, *NetworkStatusRequest) (*NetworkStatusResponse, error)
}
```

## Quick Start
This section explains an example deployment using netutil in sidecar container.

### Build

1. Compile executable:
```
$ make
```

This builds two binaries under `bin/` dir, one is named `server` that will run in sidecar container, the other is client which is an example of application that queries API server via gRPC call.

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

I0621 09:54:52.808516       6 main.go:51]  starting netutil server at: /etc/netutil/net.sock
I0621 09:54:52.811816       6 main.go:72]  netutil server start serving
I0621 09:55:04.014221       6 main.go:92]  getting cpuset from path: /proc/42/root/sys/fs/cgroup/cpuset/cpuset.cpus
I0621 09:55:04.015374       6 main.go:103] getting environment variables from path: /proc/42/environ
I0621 09:55:04.016163       6 main.go:128] getting network status from path: /etc/podinfo/annotations
...
```

5. Check for logs of application container
````
$ kubectl logs testpod -c container-2

I0621 09:55:04.011925      42 client.go:23] starting netutil client
I0621 09:55:04.013711      42 client.go:35] os.Getpid() returns: 42
I0621 09:55:04.015124      42 client.go:40] GetCPUInfo Response from NetUtilServer: 1-5
| HOSTNAME                 |: testpod5
| _                        |: /usr/bin/client
I0621 09:55:04.015849      42 client.go:46] GetEnv Response from NetUtilServer:
| container                |: docker
| KUBERNETES_PORT_443_TCP_ADDR|: 10.96.0.1
| KUBERNETES_PORT_443_TCP_PORT|: 443
| KUBERNETES_SERVICE_HOST  |: 10.96.0.1
| PWD                      |: /
| KUBERNETES_PORT_443_TCP_PROTO|: tcp
| PATH                     |: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
| SHLVL                    |: 1
| KUBERNETES_SERVICE_PORT_HTTPS|: 443
| KUBERNETES_SERVICE_PORT  |: 443
| PCIDEVICE_INTEL_COM_SRIOV|: 0000:1a:02.0,0000:1a:02.1
| KUBERNETES_PORT_443_TCP  |: tcp://10.96.0.1:443
| KUBERNETES_PORT          |: tcp://10.96.0.1:443
| INSTALL_PKGS             |: golang
| HOME                     |: /root
I0621 09:55:04.016549      42 client.go:55] GetEnv Response from NetUtilServer:
| 0                        |: ips:"10.96.1.163"
| 1                        |: name:"sriov-network" interface:"net1" ips:"10.56.217.96" mac:"da:18:1d:eb:ef:f2"
| 2                        |: name:"sriov-network" interface:"net2" ips:"10.56.217.97" mac:"fa:32:38:0e:b5:94"

````
