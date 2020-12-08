# Sample GO Application that calls Net Utility GO APIs

## Overview
This directory contains a GO program (app_sample.go) that calls the 
GO Net Utility Library. It demonstrates how the APIs can be called.
It also prints out the returned data.

## Quick Start
This section explains an example of building the GO sample application
that uses Network Utility Library.

### Compile executable
To compile the sample GO application:
```
$ cd $GOPATH/src/github.com/openshift/app-netutil/
$ make
```

This builds the GO sample application binary called `go_app` under the
`bin/` directory.

### Test
Run the application binary:
```
$ ./bin/go_app
```

This application run a forever loop, calling the Network Utility Library
and printing the results. Then sleeping for 1 minute and repeating. Use
\<CTRL\>-C to exit.

#### Debug Logs
To see additional debug messages, pass logging information as input to the
sample GO application:
```
$ ./bin/go_app -stderrthreshold=INFO
```

Valid log levels are:
* ERROR
* WARNING
* INFO

#### Run locally
If the application is not actually running in a container where annotations have been
exposed, run the following to copy a sample annotation file onto the system. There are
a couple of examples, so choose one that suits your testing. Make sure to name the
file `annotations` in the `/etc/podnetinfo/` directory.
```
$ sudo mkdir -p /etc/podnetinfo/
$ sudo cp samples/annotations/annotations_deviceinfo_pci /etc/podnetinfo/annotations
```

SR-IOV exposes the PCI Addresses of the VF to the container using an
environmental variable. If the application is not actually running in a
container where the SR-IOV environmental variables have been created, pass
them in through the command line. If the `annotations` file is using the
`device-info` field in the `network-status`, then make sure the PCI values
match.
```
$ PCIDEVICE_INTEL_COM_SRIOV=0000:01:02.5,0000:01:0a.4 ./bin/go_app -stderrthreshold=INFO
```

#### Hugepage Requests and Limits
To test the hugepage request and limit are being provided to a container via
the Downward API, the values need to be provided in the associated files.
If the application is not actually running in a container, then the files
can be created manually.

To simulate values being set by the
[SR-IOV Network-Resource-Injector](https://github.com/k8snetworkplumbingwg/network-resources-injector),
include the hugepage size and container name:
`hugepages_{1G|2M}_{request|limit}_<ContainerName>`

For example:
```
sudo sh -c 'echo "1024" >>/etc/podnetinfo/hugepages_1G_request_sriov-example'
sudo sh -c 'echo "1024" >>/etc/podnetinfo/hugepages_1G_limit_sriov-example'
```

SR-IOV Network-Resource-Injector also injects an environment variable into
the pod spec with the container's name. This allows the application to process
hugepage data properly, especially if more than one container in the pod has
requested hugepages. This can also be passed in via the command line when
running locally:
```
$ CONTAINER_NAME=sriov-example PCIDEVICE_INTEL_COM_SRIOV=0000:01:02.5,0000:01:0a.4 ./bin/go_app -stderrthreshold=INFO
```


Upstream Kubernetes hugepage Downward API examples use a simpler naming
convention (container name is still accepted): `hugepages_{request|limit}`

For example:
```
sudo sh -c 'echo "1024" >>/etc/podnetinfo/hugepages_request'
sudo sh -c 'echo "1024" >>/etc/podnetinfo/hugepages_limit'
```

app-netutil will handle both formats. Just include the environment variable
`CONTAINER_NAME` when the file names include the container name.

### Clean up
To cleanup all generated files, run:
```
$ make clean
```

This cleans up built binary and other generated files.
