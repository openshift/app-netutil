// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2021 Red Hat, Inc.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libnetutil_api.h"
#include "c_util.h"


extern void dumpInterfaces(struct InterfaceResponse *pIfaceRsp) {
	int i, j;
	bool printReturn;

	if ((pIfaceRsp) && (pIfaceRsp->pIface)) {
		for (i = 0; i < pIfaceRsp->numIfacePopulated; i++) {
			printf("  Interface[%d]:\n", i);

			printf("  ");
			printf("  DeviceType=%s",
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_HOST) ? "host" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_SRIOV) ? "SR-IOV" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_PCI) ? "PCI" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_VHOST) ? "vHost" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_MEMIF) ? "memif" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_VDPA) ? "vDPA" :
				(pIfaceRsp->pIface[i].DeviceType == NETUTIL_TYPE_UNKNOWN) ? "unknown" : "error");

			if (pIfaceRsp->pIface[i].NetworkStatus.Name) {
				printf("  Name=\"%s\"", pIfaceRsp->pIface[i].NetworkStatus.Name);
			}
			if (pIfaceRsp->pIface[i].NetworkStatus.Interface) {
				printf("  Interface=\"%s\"", pIfaceRsp->pIface[i].NetworkStatus.Interface);
			}
			printf("\n");

			printReturn = false;
			if (pIfaceRsp->pIface[i].NetworkStatus.Mac) {
				if (printReturn == false) {
					printReturn = true;
					printf("  ");
				}
				printf("  MAC=\"%s\"", pIfaceRsp->pIface[i].NetworkStatus.Mac);
			}
			for (j = 0; j < NETUTIL_NUM_IPS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.IPs[j]) {
					if (printReturn == false) {
						printReturn = true;
						printf("    DNS Nameservers: ");
					}
					printf("  IP=\"%s\"", pIfaceRsp->pIface[i].NetworkStatus.IPs[j]);
				}
			}
			if (printReturn) {
				printf("\n");
			}

			printReturn = false;
			for (j = 0; j < NETUTIL_NUM_DNS_NAMESERVERS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Nameservers[j]) {
					if (printReturn == false) {
						printReturn = true;
						printf("    DNS Nameservers: ");
					}
					printf(" \"%s\"", pIfaceRsp->pIface[i].NetworkStatus.DNS.Nameservers[j]);
				}
			}
			if (printReturn) {
				printf("\n");
			}

			if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Domain) {
				printf("    DNS Domain: \"%s\"\n", pIfaceRsp->pIface[i].NetworkStatus.DNS.Domain);
			}

			printReturn = false;
			for (j = 0; j < NETUTIL_NUM_DNS_SEARCH; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Search[j]) {
					if (printReturn == false) {
						printReturn = true;
						printf("    DNS Search: ");
					}
					printf(" \"%s\"", pIfaceRsp->pIface[i].NetworkStatus.DNS.Search[j]);
				}
			}
			if (printReturn) {
				printf("\n");
			}

			printReturn = false;
			for (j = 0; j < NETUTIL_NUM_DNS_OPTIONS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Options[j]) {
					if (printReturn == false) {
						printReturn = true;
						printf("    DNS Options: ");
					}
					printf(" \"%s\"", pIfaceRsp->pIface[i].NetworkStatus.DNS.Options[j]);
				}
			}
			if (printReturn) {
				printf("\n");
			}

			switch (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Type) {
				case NETUTIL_TYPE_PCI:
					printf("    Type=PCI");
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PciAddress) {
						printf("  PCIAddress=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PciAddress);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.Vhostnet) {
						printf("  Vhostnet=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.Vhostnet);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.RdmaDevice) {
						printf("  RdmaDevice=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.RdmaDevice);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PfPciAddress) {
						printf("  PF-PCIAddress=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PfPciAddress);
					}
					printf("\n");
					break;
				case NETUTIL_TYPE_VHOST:
					printf("    Type=vHOST");
					printf("  Mode=%s",
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Mode == NETUTIL_VHOST_MODE_CLIENT) ? "client" :
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Mode == NETUTIL_VHOST_MODE_SERVER) ? "server" : "error");
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path) {
						printf("  Path=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path);
					}
					printf("\n");
					break;
				case NETUTIL_TYPE_MEMIF:
					printf("    Type=Memif");
					printf("  Role=%s",
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Role == NETUTIL_MEMIF_ROLE_MASTER) ? "master" :
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Role == NETUTIL_MEMIF_ROLE_SLAVE) ? "slave" : "error");
					printf("  Mode=%s",
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Mode == NETUTIL_MEMIF_MODE_ETHERNET) ? "ethernet" :
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Mode == NETUTIL_MEMIF_MODE_IP) ? "ip" :
						(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Mode == NETUTIL_MEMIF_MODE_INJECT_PUNT) ? "inject-punt" : "error");
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Path) {
						printf("  Path=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Path);
					}
					printf("\n");
					break;
				case NETUTIL_TYPE_VDPA:
					printf("    Type=vDPA");
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.ParentDevice) {
						printf("  ParentDevice=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.ParentDevice);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Driver) {
						printf("  Driver=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Driver);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Path) {
						printf("  Path=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Path);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PciAddress) {
						printf("  PCIAddress=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PciAddress);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PfPciAddress) {
						printf("  PF-PCIAddress=%s", pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PfPciAddress);
					}
					printf("\n");
					break;
			}
		}
	}
}

extern void freeInterfaces(struct InterfaceResponse *pIfaceRsp) {
	int i, j;

	if ((pIfaceRsp) && (pIfaceRsp->pIface)) {
		for (i = 0; i < pIfaceRsp->numIfacePopulated; i++) {
			if (pIfaceRsp->pIface[i].NetworkStatus.Name) {
				free(pIfaceRsp->pIface[i].NetworkStatus.Name);
			}
			if (pIfaceRsp->pIface[i].NetworkStatus.Interface) {
				free(pIfaceRsp->pIface[i].NetworkStatus.Interface);
			}

			if (pIfaceRsp->pIface[i].NetworkStatus.Mac) {
				free(pIfaceRsp->pIface[i].NetworkStatus.Mac);
			}
			for (j = 0; j < NETUTIL_NUM_IPS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.IPs[j]) {
					free(pIfaceRsp->pIface[i].NetworkStatus.IPs[j]);
				}
			}

			for (j = 0; j < NETUTIL_NUM_DNS_NAMESERVERS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Nameservers[j]) {
					free(pIfaceRsp->pIface[i].NetworkStatus.DNS.Nameservers[j]);
				}
			}
			if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Domain) {
				free(pIfaceRsp->pIface[i].NetworkStatus.DNS.Domain);
			}
			for (j = 0; j < NETUTIL_NUM_DNS_SEARCH; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Search[j]) {
					free(pIfaceRsp->pIface[i].NetworkStatus.DNS.Search[j]);
				}
			}
			for (j = 0; j < NETUTIL_NUM_DNS_OPTIONS; j++) {
				if (pIfaceRsp->pIface[i].NetworkStatus.DNS.Options[j]) {
					free(pIfaceRsp->pIface[i].NetworkStatus.DNS.Options[j]);
				}
			}

			switch (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Type) {
				case NETUTIL_TYPE_PCI:
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PciAddress) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PciAddress);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.Vhostnet) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.Vhostnet);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.RdmaDevice) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.RdmaDevice);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PfPciAddress) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Pci.PfPciAddress);
					}
					break;
				case NETUTIL_TYPE_VHOST:
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.VhostUser.Path);
					}
					break;
				case NETUTIL_TYPE_MEMIF:
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Path) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Memif.Path);
					}
					break;
				case NETUTIL_TYPE_VDPA:
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.ParentDevice) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.ParentDevice);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Driver) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Driver);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Path) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.Path);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PciAddress) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PciAddress);
					}
					if (pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PfPciAddress) {
						free(pIfaceRsp->pIface[i].NetworkStatus.DeviceInfo.Vdpa.PfPciAddress);
					}
					break;
			}
		}

		free(pIfaceRsp->pIface);
	}
}

extern void dumpHugepages(struct HugepagesResponse *pHugepagesRsp) {
	int i;

	if (pHugepagesRsp) {
		if (pHugepagesRsp->MyContainerName) {
			printf("  MyContainerName=%s\n", pHugepagesRsp->MyContainerName);
		}
		if (pHugepagesRsp->pHugepages) {
			for (i = 0; i < pHugepagesRsp->numStructPopulated; i++) {
				printf("  Hugepages[%d]:\n", i);

				printf("  ");
				if (pHugepagesRsp->pHugepages[i].ContainerName) {
					printf("  ContainerName=%s", pHugepagesRsp->pHugepages[i].ContainerName);
				}
				printf("  Request: 1G=%ld 2M=%ld Ukn=%ld  Limit: 1G=%ld 2M=%ld Ukn=%ld\n",
					pHugepagesRsp->pHugepages[i].Request1G,
					pHugepagesRsp->pHugepages[i].Request2M,
					pHugepagesRsp->pHugepages[i].Request,
					pHugepagesRsp->pHugepages[i].Limit1G,
					pHugepagesRsp->pHugepages[i].Limit2M,
					pHugepagesRsp->pHugepages[i].Limit);
			}
		}
	}
}

extern void freeHugepages(struct HugepagesResponse *pHugepagesRsp) {
	int i;

	if (pHugepagesRsp) {
		if (pHugepagesRsp->MyContainerName) {
			free(pHugepagesRsp->MyContainerName);
		}
		if (pHugepagesRsp->pHugepages) {
			for (i = 0; i < pHugepagesRsp->numStructPopulated; i++) {
				if (pHugepagesRsp->pHugepages[i].ContainerName) {
					free(pHugepagesRsp->pHugepages[i].ContainerName);
				}
			}
			free(pHugepagesRsp->pHugepages);
		}
	}
}

