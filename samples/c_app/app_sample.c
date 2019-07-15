#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libnetutil_api.h"

int main() {
	struct CPUResponse cpuRsp;
	struct EnvResponse envRsp;
	struct NetworkStatusResponse networkStatusRsp;
	int i, j;

	printf("Starting sample C application.\n");


	//
	// Example of a C call to GO that returns a string.
	//
	// Note1: Calling C function must free the string.
	//
	printf("Call NetUtil GetCPUInfo():\n");
	memset(&cpuRsp, 0, sizeof(cpuRsp));
	GetCPUInfo(&cpuRsp);
	if (cpuRsp.CPUSet) {
		printf("  cpuRsp.CPUSet = %s\n", cpuRsp.CPUSet);
		free(cpuRsp.CPUSet);
	}


	//
	// Example of a C call to GO that returns a string.
	//
	// Note1: Calling C function must free the string.
	//
	printf("Call NetUtil GetEnv():\n");
	memset(&envRsp, 0, sizeof(envRsp));
	GetEnv(&envRsp);
	for (i = 0; i < NETUTIL_NUM_ENVS; i++) {
		if (envRsp.Envs[i].Index) {
			printf("  envRsp.Envs[%d].Index = %s\n", i, envRsp.Envs[i].Index);
			free(envRsp.Envs[i].Index);
		}
		if (envRsp.Envs[i].Value) {
			printf("  envRsp.Envs[%d].Value = %s\n", i, envRsp.Envs[i].Value);
			free(envRsp.Envs[i].Value);
		}
	}


	//
	// Example of a C call to GO that returns a structure
	// containing a slice of strucures which contain strings
	// and slices of strings.
	//
	// Note1: Calling C function must free the string.
	// Note2: Haven't investigated slices yet. So they
	//   defined as arrays.
	// Note3: The GO side cannot return any allocated
	//   data, so the data is allocated on the C side and
	//   passed in as a pointer.
	//
	printf("Call NetUtil GetNetworkStatus():\n");
	memset(&networkStatusRsp, 0, sizeof(networkStatusRsp));
	GetNetworkStatus(&networkStatusRsp);
	for (i = 0; i < NETUTIL_NUM_NETWORKSTATUS; i++) {
		if (networkStatusRsp.Status[i].Name) {
			printf("  networkStatusRsp.Status[%d].Name = %s\n", i, networkStatusRsp.Status[i].Name);
			free(networkStatusRsp.Status[i].Name);
		}

		if (networkStatusRsp.Status[i].Interface) {
			printf("  networkStatusRsp.Status[%d].Interface = %s\n", i, networkStatusRsp.Status[i].Interface);
			free(networkStatusRsp.Status[i].Interface);
		}

		for (j = 0; j < NETUTIL_NUM_IPS; j++) {
			if (networkStatusRsp.Status[i].IPs[j]) {
				printf("  networkStatusRsp.Status[%d].IPs[%d] = %s\n", i, j, networkStatusRsp.Status[i].IPs[j]);
				free(networkStatusRsp.Status[i].IPs[j]);
			}
		}

		if (networkStatusRsp.Status[i].Mac) {
			printf("  networkStatusRsp.Status[%d].Mac = %s\n", i, networkStatusRsp.Status[i].Mac);
			free(networkStatusRsp.Status[i].Mac);
		}
	}
}
