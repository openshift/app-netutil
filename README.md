# Network utility Library for application running on Kubernetes

## Table of Contents

- [Network Utility](#network-utility)
- [Quick Start](#quick-start)
	- [Build](#build)
	- [Create Test Pod](#create-test-pod) 

## Network Utility
Network Utility is a library that provides API methods for applications to query network information associated with pod. Network Utility is written in golang and can be built into application binary.

Currently there are three API methods implemented:

```
GetCPUInfo() (*CPUResponse, error)
GetEnv() (*EnvResponse, error)
GetNetworkStatus() (*NetworkStatusResponse, error)
```

## Quick Start
This section explains an example of building sample appication that uses Network Utility Library.

### Build

1. Compile executable:
```
$ make
```

This builds application binary called `samples` under `bin/` dir.

2. Build application container image:
```
$ make image
```

This will build an image containing application binary `samples`, the image will be used when creating test pod.

3. Clean up:
```
$ make clean
```

This cleans up built binary and softlinks.

### Create Test Pod

 1. Create application pod:
```
$ kubectl create -f deployments/pod.yaml
```
 2. Check for pod logs:
```
$ kubectl logs testpod

I0710 08:07:16.902139       1 app_sample.go:14] starting sample application
I0710 08:07:16.903046       1 resource.go:21] getting cpuset from path: /proc/1/root/sys/fs/cgroup/cpuset/cpuset.cpus
I0710 08:07:16.903574       1 app_sample.go:21] netlib.GetCPUInfo Response: 0-35
I0710 08:07:16.903599       1 resource.go:32] getting environment variables from path: /proc/1/environ
I0710 08:07:16.903669       1 app_sample.go:27] netlib.GetEnv Response:
| KUBERNETES_PORT_443_TCP_PROTO|: tcp
| KUBERNETES_PORT_443_TCP_PORT|: 443
| INSTALL_PKGS             |: golang
| PATH                     |: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
| HOSTNAME                 |: testpod
| KUBERNETES_PORT_443_TCP  |: tcp://10.96.0.1:443
| KUBERNETES_SERVICE_HOST  |: 10.96.0.1
| KUBERNETES_SERVICE_PORT  |: 443
| container                |: docker
| HOME                     |: /root
| KUBERNETES_SERVICE_PORT_HTTPS|: 443
| KUBERNETES_PORT          |: tcp://10.96.0.1:443
| KUBERNETES_PORT_443_TCP_ADDR|: 10.96.0.1
I0710 08:07:16.903756       1 network.go:18] getting network status from path: /etc/podinfo/annotations
I0710 08:07:16.904215       1 app_sample.go:36] netlib.GetNetworkStatus Response:
| 0                        |: &{Name: Interface: IPs:[10.96.1.157] Mac:}
...
```
