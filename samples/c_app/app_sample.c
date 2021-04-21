// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2021 Red Hat, Inc.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libnetutil_api.h"
#include "c_util.h"

int main() {
	struct AppNetutilConfig config;
	struct CPUResponse cpuRsp;
	struct HugepagesResponse hugepagesRsp;
	struct InterfaceResponse ifaceRsp;
	int err;

	printf("Starting sample C application.\n");

	//
	// Set APP NetUtils Settings
	//
	printf("Call NetUtil SetConfig():\n");
	memset(&config, 0, sizeof(config));

	// Update Downward API Settings
	//config.downwardAPI.baseDir = "/etc/podinfo";

	// Update Log Settings
	config.log.level = LOG_LEVEL_INFO;
	config.log.stderr = true;
	config.log.filename = "/tmp/appnetutil.log";

	err = SetConfig(&config);
	if (err) {
		printf("Couldn't set Log Settings, err code: %d\n", err);
	} else {
		printf("  Done.\n");		
	}


	//
	// Example of a C call to GO that returns a string.
	//
	// Note1: Calling C function must free the string.
	//
	printf("Call NetUtil GetCPUInfo():\n");
	memset(&cpuRsp, 0, sizeof(cpuRsp));
	err = GetCPUInfo(&cpuRsp);
	if (err) {
		printf("Couldn't get CPU info, err code: %d\n", err);
		return err;
	}
	if (cpuRsp.CPUSet) {
		printf("  cpuRsp.CPUSet = %s\n", cpuRsp.CPUSet);

		// Free the string
		free(cpuRsp.CPUSet);
	}

	//
	// Example of a C call to GO that returns a structure.
	//
	printf("Call NetUtil GetHugepages():\n");
	memset(&hugepagesRsp, 0, sizeof(struct HugepagesResponse));
	hugepagesRsp.numStructAllocated = NETUTIL_NUM_HUGEPAGES_DATA;
	hugepagesRsp.pHugepages = malloc(hugepagesRsp.numStructAllocated * sizeof(struct HugepagesData));
	if (hugepagesRsp.pHugepages) {
		memset(hugepagesRsp.pHugepages, 0, (hugepagesRsp.numStructAllocated * sizeof(struct HugepagesData)));
		err = GetHugepages(&hugepagesRsp);
		if (err) {
			printf("Couldn't get hugepage data, err code: %d\n", err);

			if (err == NETUTIL_ERRNO_SIZE_ERROR) {
				// One of the arrays wasn't sized correctly, but some data was allocated.
				// Free what was allocated.
				freeHugepages(&hugepagesRsp);
			}
			// Don't return on error. Common for hugepage data to not be available.
		} else {
			dumpHugepages(&hugepagesRsp);
			freeHugepages(&hugepagesRsp);
		}
	}

	//
	// Example of a C call to GO that returns a structure
	// containing a slice of structures which contains strings.
	//
	// Note1: Calling C function must free the string.
	// Note2: The GO side cannot return any allocated
	//   data, so the data is allocated on the C side and
	//   passed in as a pointer.
	// Note3: Instead of defining the input struct with a fixed
	//   array of entries, the C Program allocates the array
	//   dynamically. For now the number of entries are hardcoded.
	//   Later, could call GO to get the number of entries. 
	//
	printf("Call NetUtil GetInterfaces():\n");
	ifaceRsp.numIfaceAllocated = NETUTIL_NUM_NETWORKINTERFACE;
	ifaceRsp.numIfacePopulated = 0;
	ifaceRsp.pIface = malloc(ifaceRsp.numIfaceAllocated * sizeof(struct InterfaceData));
	if (ifaceRsp.pIface) {
		memset(ifaceRsp.pIface, 0, (ifaceRsp.numIfaceAllocated * sizeof(struct InterfaceData)));
		err = GetInterfaces(&ifaceRsp);
		if (err) {
			printf("Couldn't get network interface, err code: %d\n", err);

			if (err == NETUTIL_ERRNO_SIZE_ERROR) {
				// One of the arrays wasn't sized correctly, but some data was allocated.
				// Free what was allocated.
				freeInterfaces(&ifaceRsp);
			}
			return err;
		}

		dumpInterfaces(&ifaceRsp);
		freeInterfaces(&ifaceRsp);
	}

	return 0;
}
