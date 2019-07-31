#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libnetutil_api.h"

#define NETUTIL_NUM_DEFAULT_ENVS  200

int main() {
	struct CPUResponse cpuRsp;
	struct EnvResponse envRsp;
	struct NetworkStatusResponse networkStatusRsp;
	struct NetworkInterfaceResponse networkInterfaceRsp;
	int i, j;
	int num_envs;
	int err;
	char* int_type = "pci";

	printf("Starting sample C application.\n");


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
		free(cpuRsp.CPUSet);
	}


	//
	// Example of a C call to GO that returns a string.
	//
	// Note1: Calling C function must free the string.
	// Note2: Instead of defining the input struct with a fixed
	//   array of entries, the C Program allocates the array
	//   dynamically. For now hardcoded. Later, could call GO to get
	//   the number of entries. 
	//
	printf("Call NetUtil GetEnv():\n");
	num_envs = 100;
	envRsp.netutil_num_envs = num_envs;
	envRsp.pEnvs = malloc(num_envs * sizeof(struct EnvData));
	if (envRsp.pEnvs) {
		memset(envRsp.pEnvs, 0, (num_envs * sizeof(struct EnvData)));
		err = GetEnv(&envRsp);
		if (err) {
			printf("Couldn't get Env, err code: %d\n", err);
			return err;
		}
		for (i = 0; i < num_envs; i++) {
			if (envRsp.pEnvs[i].Index) {
				printf("  envRsp.pEnvs[%d].Index = %s\n", i, envRsp.pEnvs[i].Index);
				free(envRsp.pEnvs[i].Index);
			}
			if (envRsp.pEnvs[i].Value) {
				printf("  envRsp.pEnvs[%d].Value = %s\n", i, envRsp.pEnvs[i].Value);
				free(envRsp.pEnvs[i].Value);
			}
		}
		free(envRsp.pEnvs);
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
	err = GetNetworkStatus(&networkStatusRsp);
	if (err) {
		printf("Couldn't get network status, err code: %d\n", err);
		return err;
	}
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


	//
	// Example of a C call to GO that returns a structure
	// containing a slice of structures which contain strings.
	//
	// Note1: Calling C function must free the string.
	// Note2: Haven't investigated slices yet. So they
	//   defined as arrays.
	// Note3: The GO side cannot return any allocated
	//   data, so the data is allocated on the C side and
	//   passed in as a pointer.
	//
	printf("Call NetUtil GetNetworkInterface():\n");
	memset(&networkInterfaceRsp, 0, sizeof(networkInterfaceRsp));
	err = GetNetworkInterface(int_type, &networkInterfaceRsp);
	if (err) {
		printf("Couldn't get network interface, err code: %d\n", err);
		return err;
	}
	for (i = 0; i < NETUTIL_NUM_NETWORKINTERFACE; i++) {
		if (networkInterfaceRsp.Interface[i].Type) {
			printf("  networkInterfaceRsp.Interface[%d].Type = %s\n", i, networkInterfaceRsp.Interface[i].Type);
			free(networkInterfaceRsp.Interface[i].Type);
		}

		if (networkInterfaceRsp.Interface[i].ID) {
			printf("  networkInterfaceRsp.Interface[%d].ID = %s\n", i, networkInterfaceRsp.Interface[i].ID);
			free(networkInterfaceRsp.Interface[i].ID);
		}
	}

	return 0;
}
