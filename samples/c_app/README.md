# Sample C Application that calls Net Utility C APIs

## Overview
This directory contains a C program (app_sample.c) that calls the C APIs of the 
GO Net Utility Library. It demonstrates how the APIs can be called and how to
properly free any associatied memory. It also prints out the returned data.

## Quick Start
This section explains an example of building the C sample application that uses
Network Utility Library.

1. Compile executable:
```
$ cd $GOPATH/src/github.com/openshift/app-netutil/
$ make c_sample
```

This builds the GO library as a `.so` file and a C header file under the `bin/` directory.
The sample C application then includes the C header file and the application binary called
`c_sample` is built under `bin/` directory.

2. Set LD_LIBRARY_PATH
Before testing, the application needs to know where the shared library is located. Either
copy the `.so` file to a common location (i.e. `/usr/lib/`) or set `LD_LIBRARY_PATH`:
```
$ echo $LD_LIBRARY_PATH

$ export LD_LIBRARY_PATH=$PWD/bin:$LD_LIBRARY_PATH
```

Note: Printed the original value of `LD_LIBRARY_PATH` so it can be reset later if desired.
Use the following to clear out:
```
$ unset LD_LIBRARY_PATH
```

4. Test
To set, run the application binary:
```
$ ./bin/c_sample
```

If the application is not actually running in a container where annotations have been
exposed, run the following to copy a sample annotation file onto the system:
```
$ sudo mkdir -p /etc/podinfo/
$ sudo cp samples/c_app/annotations /etc/podinfo/.
```

3. Clean up:
```
$ make clean
```

This cleans up built binary and other generated files.
