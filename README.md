# Network Utility Library for Application Running on Kubernetes

## Table of Contents

- [Network Utility](#network-utility)
- [APIs](#apis)
	- [GO APIs](#go-apis)
	- [C APIs](#c-apis)
		- [C Sample APP](#c-sample-app)
		- [DPDK Sample Image](#dpdk-sample-image)
- [Quick Start](#quick-start)
	- [Build GO APP](#build-go-app)
	- [Build C APP](#build-c-app)
	- [Create testpod Image](#create-testpod-image)
	- [Create dpdk-app-centos Image](#create-dpdk-app-centos-image)

## Network Utility
Network Utility (app-netutil) is a library that provides API methods
for applications running in a container to query network information
associated with pod. Network Utility is written in golang and can be
built into an application binary. It also has a C language binding
allowing it to be built into C applications.

To add virtio based interfaces into a DPDK based application in a
container, the DPDK application needs a unix socket file, which is shared
with the host through a VolumeMount, and a set of configuration data about
how the socketfile should be used. Currently, the Userspace CNI uses
annotations or configuration files to share the data between host and
container. SR-IOV needs to get the PCI Addresses of the VFs share with the
DPDK application. Currently it is using Environmental Variables to do this.
Once the above data is in the container, this library has been written to
abstract out where to look and how to process all data passed in.

> Note: [Userspace CNI](https://github.com/intel/userspace-cni-network-plugin) is not officially supported in OpenShift.

Additional Information:
* [CONTRIBUTING.md](CONTRIBUTING.md)

## APIs
Currently there are three API methods implemented:
* `GetCPUInfo()`
  * This function determines which CPUs are available to the container
  and returns the list to the caller.
* `GetHugepages()`
  * This function determines the amount of hugepage memory available to
  the container and returns the values to the caller.
* `GetInterfaces()`
  * This function determines the set of interfaces in the container and
  returns the list, along with the interface type and type specific data.

There is a GO and C version of each of these functions.

### GO APIs

There is a GO sample app that provides an example of how to include the
app-netutil as a library in a GO program and how to use the existing APIs:
* [go_app](samples/go_app/README.md)

### C APIs

#### C Sample APP
There is a C sample app that provides an example of how to include the
app-netutil as a library in a C program and how to use the existing APIs:
* [c_app](samples/c_app/README.md)

The gotcha with the C APIs is that data must be allocated on the C side
and passed into the GO library. So the C APIs take data buffers as input.
Then the GO library populates the input structures with the collect data.
It is then up to the C side to free any allocated data.

GO has a special handling for strings where it still allocates the memory
for the string on the C side, but it is hidden in the C.CString() library.
So strings passed from the GO library back to the calling C code must also
free strings, even though there were not explicitly malloc'd. The sample C
code shows examples.

#### DPDK Sample Image
The initial problem `app-netutil` is trying to solve is to collect initial
configuration data for a DPDK application running in a container. The DPDK
Library is written in C, so there is a sample Docker Image that leverages
the C APIs of `app-netutil` to collect the initial configuration data an
then use it to start DPDK. See:
* [dpdk_app](samples/dpdk_app/dpdk-app-centos/README.md)

## Quick Start
This section provides examples of building the sample applications that use
the Network Utility Library. This is just quick start guide, more details
can be found in the links associated with each section.

### Build GO APP

1. Compile executable:
```
$ make go_sample
```
This builds application binary called `go_app` under `bin/` dir.

2. Run:
```
$ PCIDEVICE_INTEL_COM_SRIOV=0000:01:02.5,0000:01:0a.4 ./bin/go_app
| CPU     |: 0-63
| HugePage|: Request=1024  Limit=1024
| 0       |: IfName=eth0  DeviceType=host  Network=  Default=true
|         |:   MAC=06:eb:7f:3a:48:85  IPs=[10.244.0.25]  DNS={Nameservers:[] Domain: Search:[] Options:[]}
| 1       |: IfName=net1  DeviceType=sr-iov  Network=default/sriov-net-a  Default=false
|         |:   MAC=  IPs=[]  DNS={Nameservers:[] Domain: Search:[] Options:[]}
|         |:   pci  1.0.0  PCI=0000:01:02.5  PF=  Vhostnet=  RdmaDevice=
| 2       |: IfName=net2  DeviceType=sr-iov  Network=default/sriov-net-b  Default=false
|         |:   MAC=  IPs=[]  DNS={Nameservers:[] Domain: Search:[] Options:[]}
|         |:   pci  1.0.0  PCI=0000:01:0a.4  PF=  Vhostnet=  RdmaDevice=

<CTRL>-C
```

3. Clean up:
```
$ make clean
```

This cleans up built binary and softlinks.


For more details, see:
* [go_app](samples/go_app/README.md)

### Build C APP

1. Compile executable:
```
$ make c_sample
```
This builds application binary called `c_sample` under `bin/`
directory. The `bin/` directory also contains the C header file
`libnetutil_api.h`  and shared library `libnetutil_api.so` needed
to build the C APP.

2. Run:
```
$ PCIDEVICE_INTEL_COM_SRIOV=0000:01:02.5,0000:01:0a.4 \
  LD_LIBRARY_PATH=$PWD/bin:$LD_LIBRARY_PATH \
  ./bin/c_sample
Starting sample C application.
Call NetUtil GetCPUInfo():
  cpuRsp.CPUSet = 0-63
Call NetUtil GetHugepages():
  Request = 1024  Limit = 1024
Call NetUtil GetInterfaces():
  Interface[0]:
    DeviceType=host  Interface="eth0"
    MAC="06:eb:7f:3a:48:85"  IP="10.244.0.25"
  Interface[1]:
    DeviceType=SR-IOV  Name="default/sriov-net-a"  Interface="net1"
    Type=PCI  PCIAddress=0000:01:02.5
  Interface[2]:
    DeviceType=SR-IOV  Name="default/sriov-net-b"  Interface="net2"
    Type=PCI  PCIAddress=0000:01:0a.4
```

3. Clean up:
```
$ make clean
```

This cleans up built binary and softlinks.


For more details, see:
* [c_app](samples/c_app/README.md)


### Create testpod Image
The `testpod` image is a CentOS base image built the `app-netutil`
library. It simply creates a container that runs the `go_app` sample
applicatation described above.

1. Build application container image:
```
$ make testpod
```
2. Create application pod:
```
$ kubectl create -f samples/testpod/pod.yaml
```
3. Check for pod logs:
```
$ kubectl logs testpod
I1202 18:47:09.067139       1 app_sample.go:16] starting sample application
I1202 18:47:09.069655       1 app_sample.go:20] CALL netlib.GetCPUInfo:
I1202 18:47:09.070542       1 resource.go:27] getting cpuset from path: /proc/1/root/sys/fs/cgroup/cpuset/cpuset.cpus
| CPU     |: 0-63
I1202 18:47:09.071110       1 app_sample.go:26] netlib.GetCPUInfo Response:
I1202 18:47:09.071126       1 app_sample.go:30] CALL netlib.GetHugepages:
I1202 18:47:09.071149       1 hugepages.go:22] GetHugepages: Open /etc/podnetinfo/hugepages_request
I1202 18:47:09.071349       1 hugepages.go:25] Error getting /etc/podnetinfo/hugepages_request info: open /etc/podnetinfo/hugepages_request: no such file or directory
I1202 18:47:09.071363       1 hugepages.go:35] GetHugepages: Open /etc/podnetinfo/hugepages_limit
I1202 18:47:09.071390       1 hugepages.go:38] Error getting /etc/podnetinfo/hugepages_limit info: open /etc/podnetinfo/hugepages_limit: no such file or directory
I1202 18:47:09.071404       1 app_sample.go:33] Error calling netlib.GetHugepages: open /etc/podnetinfo/hugepages_request: no such file or directory
I1202 18:47:09.071418       1 app_sample.go:40] CALL netlib.GetInterfaces:
I1202 18:47:09.071424       1 network.go:39] GetInterfaces: ENTER
I1202 18:47:09.071433       1 network.go:45] GetInterfaces: Open /etc/podnetinfo/annotations
I1202 18:47:09.071821       1 network.go:68]   s-k8s.v1.cni.cncf.io/network-status="[{\n    \"name\": \"\",\n    \"interface\": \"eth0\",\n    \"ips\": [\n        \"10.244.0.57\"\n    ],\n    \"mac\": \"82:33:bb:68:3a:3c\",\n    \"default\": true,\n    \"dns\": {}\n}]"
I1202 18:47:09.071832       1 network.go:72]   PartsLen-2
I1202 18:47:09.071843       1 network.go:74]   parts[0]-k8s.v1.cni.cncf.io/network-status
I1202 18:47:09.072035       1 network.go:68]   s-k8s.v1.cni.cncf.io/networks-status="[{\n    \"name\": \"\",\n    \"interface\": \"eth0\",\n    \"ips\": [\n        \"10.244.0.57\"\n    ],\n    \"mac\": \"82:33:bb:68:3a:3c\",\n    \"default\": true,\n    \"dns\": {}\n}]"
I1202 18:47:09.072047       1 network.go:72]   PartsLen-2
I1202 18:47:09.072054       1 network.go:74]   parts[0]-k8s.v1.cni.cncf.io/networks-status
I1202 18:47:09.072067       1 network.go:68]   s-kubernetes.io/config.seen="2020-12-02T13:46:08.010504651-05:00"
I1202 18:47:09.072075       1 network.go:72]   PartsLen-2
I1202 18:47:09.072082       1 network.go:74]   parts[0]-kubernetes.io/config.seen
I1202 18:47:09.072090       1 network.go:68]   s-kubernetes.io/config.source="api"
I1202 18:47:09.072097       1 network.go:72]   PartsLen-2
I1202 18:47:09.072105       1 network.go:74]   parts[0]-kubernetes.io/config.source
I1202 18:47:09.072111       1 networkstatus.go:42] PRINT EACH NetworkStatus - len=1
I1202 18:47:09.072118       1 networkstatus.go:46]   status:
I1202 18:47:09.072126       1 networkstatus.go:47] { eth0 [10.244.0.57] 82:33:bb:68:3a:3c true {[]  [] []} <nil>}
I1202 18:47:09.072162       1 userspace.go:49] PRINT EACH Userspace MappedDir
I1202 18:47:09.072166       1 userspace.go:50]   usrspMappedDir:
I1202 18:47:09.072171       1 userspace.go:51] 
I1202 18:47:09.072176       1 userspace.go:53] PRINT EACH Userspace ConfigData
I1202 18:47:09.072183       1 network.go:112] PROCESS ENV:
I1202 18:47:09.072190       1 resource.go:38] getting environment variables from path: /proc/1/environ
I1202 18:47:09.072232       1 network.go:165] eth0 is the "default" interface, mark as "host"
I1202 18:47:09.072246       1 network.go:221] RESPONSE:
I1202 18:47:09.072258       1 network.go:223] &{host { eth0 [10.244.0.57] 82:33:bb:68:3a:3c true {[]  [] []} <nil>}}
I1202 18:47:09.072280       1 app_sample.go:46] netlib.GetInterfaces Response:
| 0       |: IfName=eth0  DeviceType=host  Network=  Default=true
|         |:   MAC=82:33:bb:68:3a:3c  IPs=[10.244.0.57]  DNS={Nameservers:[] Domain: Search:[] Options:[]}
...
```

> NOTE: If the hugepage Downward API is not included in the Pod Spec, which
is the case for [samples/testpod/pod.yaml](samples/testpod/pod.yaml), then
the hugepage annotation files will not exist and hugepage data will not be
retrieved. The hugepage Downward API requires Kubernetes 1.20 or greater and
the feature gate to be enable.

4. Delete application pod:
```
$ kubectl delete -f deployments/pod.yaml
```


For more details, see:
* [go_app](samples/go_app/README.md)

### Create dpdk-app-centos Image
The `dpdk-app-centos` image is a CentOS base image built with DPDK
and includes the `app-netutil` library. The setup to run the image
is more complicated and depends on if you are using vhost interfaces
from something like a Userspace CNI or SR-IOV VFs from SR-IOV CNI.
Below is the quick command to build the image, but it is recommended
that additional README files are consulted for detailed setup
instructions.

1. Build application container image:
```
$ make dpdk_app
```

For more details, see:
* [dpdk_app image](samples/dpdk_app/dpdk-app-centos/README.md)
* [SR-IOV VF Deployment](samples/dpdk_app/sriov/README.md)
