#  Docker Image: dpdk-app-centos
This directory contains the files needed to build a DPDK based test image.
This image is based on CentOS 7 base image built with DPDK 19.08.
This image is intended to work with multiple CNIs that intend to tie into
DPDK running in a container.

## SR-IOV CNI and SR-IOV Device Plugin
The SR-IOV Device Plugin detects and tracks the VFs associated with an
SR-IOV PF. The SR-IOV CNI updates the VF as needed and the PCI Address
associated with the VF is passed to the container via environmental
variables. The container, like this one, boots up and reads the
environmental variables (via the app-netulit) and runs a DPDK application.

## Userspace CNI
The Userspace CNI inconjunction with the OVS CNI Library (cniovs) or VPP
CNI Library (cnivpp) creates interfaces on the host, like a vhost-user or
a memif interface, adds the host side of the interface to a local network,
like a L2 bridge, then copies information needed in the container into
annotations. The container, like this one, boots up and reads the
annotatations (via the app-netutil) and runs a DPDK application.


# Build Instructions for dpdk-app-centos Docker Image
Get the **app-netutil** repo:
```
   cd $GOPATH/src/
   go get github.com/openshift/app-netutil
```

Build the docker image:
```
   cd $GOPATH/src/github.com/openshift/app-netutil/samples/dpdk_app/dpdk-app-centos/
   docker build --rm -t dpdk-app-centos .
```
OR
```
   cd $GOPATH/src/github.com/openshift/app-netutil/
   make dpdk_app
```

## Reduce Image Size
Multi-stage builds are a new feature requiring **Docker 17.05** or higher on
the daemon and client. If multi-stage builds are NOT supported on your system,
then comment out the lines between `# BEGIN` and `# END` in the Dockerfile:
```
:

# -------- Import stage.
# Docker 17.05 or higher
# BEGIN
FROM centos
COPY --from=0 /usr/bin/dpdk-app /usr/bin/dpdk-app
COPY --from=0 /usr/bin/l2fwd /usr/bin/l2fwd
COPY --from=0 /usr/bin/l3fwd /usr/bin/l3fwd
COPY --from=0 /usr/bin/testpmd /usr/bin/testpmd
COPY --from=0 /lib64/libnetutil_api.so /lib64/libnetutil_api.so
COPY --from=0 /usr/lib64/libnuma.so.1 /usr/lib64/libnuma.so.1
# END

:
```

# Docker Image Details
This Docker Image is downloading DPDK (version 19.08 to get memif PMD)
and building it. Once built, it changes into the DPDK `testpmd`
directory (`${DPDK_DIR}/app/test-pmd`) and builds it. It then repeats
for the DPDK `l2fwd` directory (`${DPDK_DIR}/examples/l2fwd`) and the
the DPDK `l3fwd` directory (`${DPDK_DIR}/examples/l3fwd`).

`testpmd`, `l2fwd` and `l3fwd` are DPDK sample applications. `testpmd`
and `l2fwd` performs Layer 2 switching and `l3fwd` performs Layer 3
routing. All three applications are built with the `app-netutil` and
copied into `/usr/bin/`. `l3fwd` is then also copied to `/usr/bin/` and
renamed as `dpdk-app` (default option).

Which DPDK sample application is run is controlled by an environmental
variable (`DPDK_SAMPLE_APP`) set in the pod spec. If not set, the image
defaults to using `dpdk-app`, which is mapped to `l3fwd`. See
`sample/dpdk_app/sriov/sriov-pod-1.yaml` for an example of the environmental
variable `DPDK_SAMPLE_APP` being used.

Typically, `l3fwd` is started with a set of input parameters that
initializes DPDK.

For example:
```
$ l3fwd -n 4 -l 1 --master-lcore 1 -w 0000:01:0a.6 -w 0000:01:02.1 -- -p 0x3 -P --config="(0,0,1),(1,0,1)" --parse-ptype
```

This Docker image is tweaking this a little. Before `l3fwd` is built, the
main.c file (contains `main()`) is updated using `sed`. See
`l3fwd_substitute.sh`.

**NOTE:** If a different version of DPDK is needed or used, this script and
text file may need to be synchronized with the updated version. 

An additional file, dpdk-args.c, is also added to the directory and Makefile.
The changes to main.c are simply to call a function in dpdk-args.c which
will generate this list of input parameters, and then pass this private set
of parameters to DPDK functions instead of the inpupt `argc` and `argv`. When
the generated binary is copied to `/usr/bin/`, it is renamed to `dpdk-app`.

The code is leveraging this project, app-netutil
(https://github.com/openshift/app-netutil), which is a library designed to be
called within a container to collect all the configuration data, like that
stored in environmental variables by SR-IOV Device Plugin and annotations by
Userspace CNI, and expose it to a DPDK application in a clean API.

**NOTE:** For debugging, if `dpdk-app` is called with a set of input parameters,
it will skip the dpdk-args.c code and behave exactly as `l3fwd`. Just add
the `sleep` to the pod spec:
```
:
    resources:
      requests:
        memory: 2Mi
      limits:
        hugepages-2Mi: 1024Mi
    command: ["sleep", "infinity"]    <-- UNCOMMENT
  volumes:
:
```

Then get a pod shell:
```
   kubectl exec -it sriov-pod-1 -- sh
```

Run `dpdk-app` with no parameters, and it will be as if it is called
as the container is started. It also prints out the generated parameter
list, which include the dynamic socketfile path:
```
sh-4.2# dpdk-app
ENTER dpdk-app:
 argc=1
 dpdk-app
  cpuRsp.CPUSet = 0-63
  Hugepage: Request = 2048 Limit = 2048  Using = 1024
  Interface[0]:
    DeviceType=host  Interface="eth0"
    MAC="66:34:b2:6b:84:e0"  IP="10.244.0.52"
  Interface[1]:
    DeviceType=SR-IOV  Name="default/sriov-net-a"  Interface="net1"
    Type=PCI  PCIAddress=0000:01:0a.1
  Interface[2]:
    DeviceType=SR-IOV  Name="default/sriov-net-b"  Interface="net2"
    Type=PCI  PCIAddress=0000:01:02.2
 myArgc=17
 dpdk-app -m 1024 -n 4 -l 1 --master-lcore 1 -w 0000:01:0a.1 -w 0000:01:02.2 -- -p 0x3 -P --config="(0,0,1),(1,0,1)" --parse-ptype
EAL: Detected 64 lcore(s)
EAL: Detected 2 NUMA nodes
EAL: Multi-process socket /var/run/dpdk/rte/mp_socket
:
```

Then `\<CTRL-C\>` to exit and re-run `dpdk-app` with input parameters
modified as needed:
```
dpdk-app -m 1024 -n 4 -l 1 --master-lcore 1 -w 0000:01:0a.1 -w 0000:01:02.2 -- -p 0x3 -P --config="(0,0,1),(1,0,1)" --parse-ptype
```

The output from running `dpdk-app`' in the container is described here. The
"argc=1" and "dpdk-app" are a reprint of how the comamnd was called from the
commandline. As mentioned above, if "argc" is greater than 1, then all the
app-netutil code is skipped and the dpdk sample (`l3fwd` or `testpmd`) are
called as the normally would be.
```
:
ENTER dpdk-app:
 argc=1
 dpdk-app
:
```

The next set of output come from `app-netutil` and indicate what
data it has collected from the environment variables and
annotations. 
```
  cpuRsp.CPUSet = 0-63
  Hugepage: Request = 2048 Limit = 2048  Using = 1024
  Interface[0]:
    DeviceType=host  Interface="eth0"
    MAC="66:34:b2:6b:84:e0"  IP="10.244.0.52"
  Interface[1]:
    DeviceType=SR-IOV  Name="default/sriov-net-a"  Interface="net1"
    Type=PCI  PCIAddress=0000:01:0a.1
  Interface[2]:
    DeviceType=SR-IOV  Name="default/sriov-net-b"  Interface="net2"
    Type=PCI  PCIAddress=0000:01:02.2
```

The next set of output indicated how `dpdk-app` is called with
the set of parameters printed. This can be copied and rerun with
updates as needed. But the dynamic data, such as PCI Address of
SR-IOV interfaces, or vhost socketfiles are printed and can be
leveraged on subsequent runs.
```
 myArgc=15
 dpdk-app -m 1024 -n 4 -l 1 --master-lcore 1 -w 0000:01:0a.1 -w 0000:01:02.2 -- -p 0x3 -P --config="(0,0,1),(1,0,1)" --parse-ptype
```

The remaining output is from DPDK.

# Deploy Image
## SR-IOV Deployment
An example of using this Docker image with SR-IOV can be found in this
same repo. See:
 * [dpdk-app-centos](../sriov/README.md)

## Userspace CNI Deployment
An example of using this Docker image with Userspace CNI can be found in
the Userspace CNI repo. See:
* [dpdk-app-centos with Userspace CNI](https://github.com/intel/userspace-cni-network-plugin/blob/master/docker/dpdk-app-centos/)
