# Sample C Application that calls Net Utility C APIs

## Overview
This directory contains a C program (app_sample.c) that calls the C APIs of the 
GO Net Utility Library. It demonstrates how the APIs can be called and how to
properly free any associatied memory. It also prints out the returned data.

## Quick Start
This section explains an example of building the C sample application that uses
Network Utility Library.

### Compile executable:
To compile the sample C application:
```
$ cd $GOPATH/src/github.com/openshift/app-netutil/
$ make c_sample
```

This builds the GO library as a `.so` file and a C header file under the `bin/`
directory. The sample C application then includes the C header file and the
application binary called `c_sample` is built under `bin/` directory.

### Test
Before testing, the application needs to know where the shared library is located.
Either copy the `.so` file to a common location (i.e. `/usr/lib/`) or pass the
`LD_LIBRARY_PATH` environment variable via the command line.

Run the application binary:
```
$ LD_LIBRARY_PATH=$PWD/bin:$LD_LIBRARY_PATH ./bin/c_sample
```

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
$ LD_LIBRARY_PATH=$PWD/bin:$LD_LIBRARY_PATH PCIDEVICE_INTEL_COM_SRIOV=0000:01:02.5,0000:01:0a.4 ./bin/c_sample
```

#### Hugepage Requests and Limits
To test the hugepage request and limit are being provided to a container via
via the Downward API, the values need to be provided in the associated files.
If the application is not actually running in a container, then the files
can be created manually:
```
sudo sh -c 'echo "1024" >>/etc/podnetinfo/hugepages_request'
sudo sh -c 'echo "1024" >>/etc/podnetinfo/hugepages_limit'
```

### Clean up
To cleanup all generated files, run:
```
$ make clean
```

This cleans up built binary and other generated files.
