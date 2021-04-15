package main

/*
#include <stdint.h>
#include <stdbool.h>

// Mapped from app-netutil/lib/v1alpha/types.go

struct CPUResponse {
	char*    CPUSet;
};

#define NETUTIL_NUM_HUGEPAGES_DATA  10

struct HugepagesData {
	char*    ContainerName;
	int64_t  Request;
	int64_t  Limit;
	int64_t  Request1G;
	int64_t  Limit1G;
	int64_t  Request2M;
	int64_t  Limit2M;
};

// *pIface is an array of 'struct HugepagesData' that is allocated
// from the C program.
struct HugepagesResponse {
	int                   numStructAllocated;
	int                   numStructPopulated;
	char*                 MyContainerName;
	struct HugepagesData *pHugepages;
};

#define NETUTIL_ERRNO_SUCCESS 0
#define NETUTIL_ERRNO_FAIL 1
#define NETUTIL_ERRNO_SIZE_ERROR 2


#define NETUTIL_NUM_IPS               10
#define NETUTIL_NUM_NETWORKINTERFACE  10
#define NETUTIL_NUM_DNS_NAMESERVERS    5
#define NETUTIL_NUM_DNS_SEARCH         5
#define NETUTIL_NUM_DNS_OPTIONS        5

struct DNS {
	char*  Nameservers[NETUTIL_NUM_DNS_NAMESERVERS];
	char*  Domain;
	char*  Search[NETUTIL_NUM_DNS_SEARCH];
	char*  Options[NETUTIL_NUM_DNS_OPTIONS];
};

struct PciDevice {
	char*  PciAddress;
	char*  Vhostnet;
	char*  RdmaDevice;
	char*  PfPciAddress;
};

struct VdpaDevice {
	char*  ParentDevice;
	char*  Driver;
	char*  Path;
	char*  PciAddress;
	char*  PfPciAddress;
};

#define NETUTIL_VHOST_MODE_CLIENT  0
#define NETUTIL_VHOST_MODE_SERVER  1
struct VhostDevice {
	int    Mode;
	char*  Path;
};

#define NETUTIL_MEMIF_ROLE_MASTER       0
#define NETUTIL_MEMIF_ROLE_SLAVE        1
#define NETUTIL_MEMIF_MODE_ETHERNET     0
#define NETUTIL_MEMIF_MODE_IP           1
#define NETUTIL_MEMIF_MODE_INJECT_PUNT  2
struct MemifDevice {
	int    Role;
	char*  Path;
	int    Mode;
};

#define NETUTIL_TYPE_UNKNOWN  0
#define NETUTIL_TYPE_HOST     1
#define NETUTIL_TYPE_SRIOV    2
#define NETUTIL_TYPE_PCI      3
#define NETUTIL_TYPE_VHOST    4
#define NETUTIL_TYPE_MEMIF    5
#define NETUTIL_TYPE_VDPA     8

struct DeviceInfo {
	int    Type;
	char*  Version;
	struct PciDevice   Pci;
	struct VdpaDevice  Vdpa;
	struct VhostDevice VhostUser;
	struct MemifDevice Memif;
};

struct NetworkStatus {
	char*  Name;
	char*  Interface;
	char*  IPs[NETUTIL_NUM_IPS];
	char*  Mac;
	int    Default;
	struct DNS DNS;
	struct DeviceInfo DeviceInfo;
};

struct InterfaceData {
	int    DeviceType;
	struct NetworkStatus NetworkStatus;
};

// *pIface is an array of 'struct InterfaceData' that is allocated
// from the C program.
struct InterfaceResponse {
	int                   numIfaceAllocated;
	int                   numIfacePopulated;
	struct InterfaceData *pIface;
};

// Mapped from app-netutil/pkg/logging/logging.go
#define LOG_LEVEL_ERROR    0
#define LOG_LEVEL_WARNING  1
#define LOG_LEVEL_INFO     2
#define LOG_LEVEL_DEBUG    3
#define LOG_LEVEL_MAX      4

struct LoggingData {
	uint32_t  stderr;
	uint32_t  level;
	char*     filename;
};

struct AppNetutilConfig {
	struct LoggingData log;
};

*/
import "C"
import "unsafe"

import (
	nettypes "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	netlib "github.com/openshift/app-netutil/lib/v1alpha"
	"github.com/openshift/app-netutil/pkg/logging"
	"github.com/openshift/app-netutil/pkg/types"
)

const (
	netutil_num_ips = 10

	// Interface type
	NETUTIL_INTERFACE_TYPE_HOST  = types.INTERFACE_TYPE_HOST
	NETUTIL_INTERFACE_TYPE_SRIOV = types.INTERFACE_TYPE_SRIOV
	NETUTIL_INTERFACE_TYPE_PCI   = nettypes.DeviceInfoTypePCI
	NETUTIL_INTERFACE_TYPE_VHOST = nettypes.DeviceInfoTypeVHostUser
	NETUTIL_INTERFACE_TYPE_MEMIF = nettypes.DeviceInfoTypeMemif
	NETUTIL_INTERFACE_TYPE_VDPA  = nettypes.DeviceInfoTypeVDPA

	// Errno
	NETUTIL_ERRNO_SUCCESS    = 0
	NETUTIL_ERRNO_FAIL       = 1
	NETUTIL_ERRNO_SIZE_ERROR = 2
)

//export GetCPUInfo
func GetCPUInfo(c_cpuResp *C.struct_CPUResponse) int64 {
	cpuRsp, err := netlib.GetCPUInfo()

	if err == nil {
		c_cpuResp.CPUSet = C.CString(cpuRsp.CPUSet)
		return NETUTIL_ERRNO_SUCCESS
	}
	_ = logging.Errorf("netlib.GetCPUInfo() err: %+v", err)
	return NETUTIL_ERRNO_FAIL
}

//export GetHugepages
func GetHugepages(c_hugepagesResp *C.struct_HugepagesResponse) int64 {

	var j C.int

	hugepagesRsp, err := netlib.GetHugepages()

	if err == nil {
		if hugepagesRsp.MyContainerName != "" {
			c_hugepagesResp.MyContainerName = C.CString(hugepagesRsp.MyContainerName)
		}

		// Map the input pointer to array of structures, c_hugepagesResp.pHugepages, to
		// a slice of the structures, c_hugepagesResp.pHugepages. Then the slice can be
		// indexed.
		c_hugepagesResp_pHugepages := (*[1 << 30]C.struct_HugepagesData)(unsafe.Pointer(c_hugepagesResp.pHugepages))[:c_hugepagesResp.numStructAllocated:c_hugepagesResp.numStructAllocated]

		for i, hugepagesData := range hugepagesRsp.Hugepages {
			if j < c_hugepagesResp.numStructAllocated {
				if hugepagesData.ContainerName != "" {
					c_hugepagesResp_pHugepages[j].ContainerName = C.CString(hugepagesData.ContainerName)
				}
				c_hugepagesResp_pHugepages[j].Request = C.long(hugepagesData.Request)
				c_hugepagesResp_pHugepages[j].Limit = C.long(hugepagesData.Limit)
				c_hugepagesResp_pHugepages[j].Request1G = C.long(hugepagesData.Request1G)
				c_hugepagesResp_pHugepages[j].Limit1G = C.long(hugepagesData.Limit1G)
				c_hugepagesResp_pHugepages[j].Request2M = C.long(hugepagesData.Request2M)
				c_hugepagesResp_pHugepages[j].Limit2M = C.long(hugepagesData.Limit2M)

				c_hugepagesResp.numStructPopulated++

				j++
			} else {

				_ = logging.Errorf("HugepagesResponse struct not sized properly."+
					"At index %d.", i)

				return NETUTIL_ERRNO_SIZE_ERROR
			}
		}
		return NETUTIL_ERRNO_SUCCESS
	}
	_ = logging.Errorf("netlib.GetHugepages() err: %+v", err)
	return NETUTIL_ERRNO_FAIL
}

//export GetInterfaces
func GetInterfaces(c_ifaceRsp *C.struct_InterfaceResponse) int64 {

	var j C.int

	ifaceRsp, err := netlib.GetInterfaces()

	if err == nil {
		j = 0

		// Map the input pointer to array of structures, c_ifaceResp.pIface, to
		// a slice of the structures, c_ifaceResp_pIface. Then the slice can be
		// indexed.
		c_ifaceResp_pIface := (*[1 << 30]C.struct_InterfaceData)(unsafe.Pointer(c_ifaceRsp.pIface))[:c_ifaceRsp.numIfaceAllocated:c_ifaceRsp.numIfaceAllocated]

		for i, iface := range ifaceRsp.Interface {
			if j < c_ifaceRsp.numIfaceAllocated {

				// Map InterfaceData
				switch iface.DeviceType {
				case NETUTIL_INTERFACE_TYPE_HOST:
					c_ifaceResp_pIface[j].DeviceType = C.NETUTIL_TYPE_HOST
				case NETUTIL_INTERFACE_TYPE_SRIOV:
					c_ifaceResp_pIface[j].DeviceType = C.NETUTIL_TYPE_SRIOV
				case NETUTIL_INTERFACE_TYPE_PCI:
					c_ifaceResp_pIface[j].DeviceType = C.NETUTIL_TYPE_PCI
				case NETUTIL_INTERFACE_TYPE_VHOST:
					c_ifaceResp_pIface[j].DeviceType = C.NETUTIL_TYPE_VHOST
				case NETUTIL_INTERFACE_TYPE_MEMIF:
					c_ifaceResp_pIface[j].DeviceType = C.NETUTIL_TYPE_MEMIF
				case NETUTIL_INTERFACE_TYPE_VDPA:
					c_ifaceResp_pIface[j].DeviceType = C.NETUTIL_TYPE_VDPA
				default:
					c_ifaceResp_pIface[j].DeviceType = C.NETUTIL_TYPE_UNKNOWN
				}

				// Map InterfaceData.NetworkStatus
				if iface.NetworkStatus.Name != "" {
					c_ifaceResp_pIface[j].NetworkStatus.Name = C.CString(iface.NetworkStatus.Name)
				}
				if iface.NetworkStatus.Interface != "" {
					c_ifaceResp_pIface[j].NetworkStatus.Interface = C.CString(iface.NetworkStatus.Interface)
				}

				for k, ip := range iface.NetworkStatus.IPs {
					if k < netutil_num_ips {
						c_ifaceResp_pIface[j].NetworkStatus.IPs[k] = C.CString(ip)
					} else {
						_ = logging.Errorf("NetworkStatus.IPs array not sized properly."+
							"At Interface %d, IP index %d.", i, k)
						return NETUTIL_ERRNO_SIZE_ERROR
					}
				}

				if iface.NetworkStatus.Mac != "" {
					c_ifaceResp_pIface[j].NetworkStatus.Mac = C.CString(iface.NetworkStatus.Mac)
				}

				if iface.NetworkStatus.Default {
					c_ifaceResp_pIface[j].NetworkStatus.Default = 1
				} else {
					c_ifaceResp_pIface[j].NetworkStatus.Default = 0
				}

				for k, nameserver := range iface.NetworkStatus.DNS.Nameservers {
					if k < C.NETUTIL_NUM_DNS_NAMESERVERS {
						c_ifaceResp_pIface[j].NetworkStatus.DNS.Nameservers[k] = C.CString(nameserver)
					} else {
						_ = logging.Errorf("NetworkStatus.DNS.Nameservers array not sized properly."+
							"At Interface %d, index %d.", i, k)
						return NETUTIL_ERRNO_SIZE_ERROR
					}
				}
				if iface.NetworkStatus.DNS.Domain != "" {
					c_ifaceResp_pIface[j].NetworkStatus.DNS.Domain = C.CString(iface.NetworkStatus.DNS.Domain)
				}
				for k, search := range iface.NetworkStatus.DNS.Search {
					if k < C.NETUTIL_NUM_DNS_SEARCH {
						c_ifaceResp_pIface[j].NetworkStatus.DNS.Search[k] = C.CString(search)
					} else {
						_ = logging.Errorf("NetworkStatus.DNS.Search array not sized properly."+
							"At Interface %d, index %d.", i, k)
						return NETUTIL_ERRNO_SIZE_ERROR
					}
				}
				for k, option := range iface.NetworkStatus.DNS.Options {
					if k < C.NETUTIL_NUM_DNS_OPTIONS {
						c_ifaceResp_pIface[j].NetworkStatus.DNS.Options[k] = C.CString(option)
					} else {
						_ = logging.Errorf("NetworkStatus.DNS.Options array not sized properly."+
							"At Interface %d, index %d.", i, k)
						return NETUTIL_ERRNO_SIZE_ERROR
					}
				}

				// Map InterfaceData.NetworkStatus.DeviceInfo
				if iface.NetworkStatus.DeviceInfo != nil {
					c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Version =
						C.CString(iface.NetworkStatus.DeviceInfo.Version)
					switch iface.NetworkStatus.DeviceInfo.Type {
					case NETUTIL_INTERFACE_TYPE_PCI:
						c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Type = C.NETUTIL_TYPE_PCI
						if iface.NetworkStatus.DeviceInfo.Pci != nil {
							if iface.NetworkStatus.DeviceInfo.Pci.PciAddress != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Pci.PciAddress =
									C.CString(iface.NetworkStatus.DeviceInfo.Pci.PciAddress)
							}
							if iface.NetworkStatus.DeviceInfo.Pci.Vhostnet != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Pci.Vhostnet =
									C.CString(iface.NetworkStatus.DeviceInfo.Pci.Vhostnet)
							}
							if iface.NetworkStatus.DeviceInfo.Pci.RdmaDevice != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Pci.RdmaDevice =
									C.CString(iface.NetworkStatus.DeviceInfo.Pci.RdmaDevice)
							}
							if iface.NetworkStatus.DeviceInfo.Pci.PfPciAddress != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Pci.PfPciAddress =
									C.CString(iface.NetworkStatus.DeviceInfo.Pci.PfPciAddress)
							}
						} else {
							_ = logging.Errorf("Error: type set to pci, but no associated DeviceInfo data")
						}
					case NETUTIL_INTERFACE_TYPE_VHOST:
						c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Type = C.NETUTIL_TYPE_VHOST
						if iface.NetworkStatus.DeviceInfo.VhostUser != nil {
							if iface.NetworkStatus.DeviceInfo.VhostUser.Mode == nettypes.VhostDeviceModeClient {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.VhostUser.Mode = C.NETUTIL_VHOST_MODE_CLIENT
							} else {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.VhostUser.Mode = C.NETUTIL_VHOST_MODE_SERVER
							}
							if iface.NetworkStatus.DeviceInfo.VhostUser.Path != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.VhostUser.Path =
									C.CString(iface.NetworkStatus.DeviceInfo.VhostUser.Path)
							}
						} else {
							_ = logging.Errorf("Error: type set to vHost, but no associated DeviceInfo data")
						}
					case NETUTIL_INTERFACE_TYPE_MEMIF:
						c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Type = C.NETUTIL_TYPE_MEMIF
						if iface.NetworkStatus.DeviceInfo.Memif != nil {
							if iface.NetworkStatus.DeviceInfo.Memif.Role == nettypes.MemifDeviceRoleMaster {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Memif.Role = C.NETUTIL_MEMIF_ROLE_MASTER
							} else {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Memif.Role = C.NETUTIL_MEMIF_ROLE_SLAVE
							}
							if iface.NetworkStatus.DeviceInfo.Memif.Mode == nettypes.MemifDeviceModeEthernet {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Memif.Mode = C.NETUTIL_MEMIF_MODE_ETHERNET
							} else if iface.NetworkStatus.DeviceInfo.Memif.Mode == nettypes.MemitDeviceModeIP {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Memif.Mode = C.NETUTIL_MEMIF_MODE_IP
							} else {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Memif.Mode = C.NETUTIL_MEMIF_MODE_INJECT_PUNT
							}
							if iface.NetworkStatus.DeviceInfo.Memif.Path != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Memif.Path =
									C.CString(iface.NetworkStatus.DeviceInfo.Memif.Path)
							}
						} else {
							_ = logging.Errorf("Error: type set to memif, but no associated DeviceInfo data")
						}
					case NETUTIL_INTERFACE_TYPE_VDPA:
						c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Type = C.NETUTIL_TYPE_VDPA
						if iface.NetworkStatus.DeviceInfo.Vdpa != nil {
							if iface.NetworkStatus.DeviceInfo.Vdpa.ParentDevice != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Vdpa.ParentDevice =
									C.CString(iface.NetworkStatus.DeviceInfo.Vdpa.ParentDevice)
							}
							if iface.NetworkStatus.DeviceInfo.Vdpa.Driver != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Vdpa.Driver =
									C.CString(iface.NetworkStatus.DeviceInfo.Vdpa.Driver)
							}
							if iface.NetworkStatus.DeviceInfo.Vdpa.Path != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Vdpa.Path =
									C.CString(iface.NetworkStatus.DeviceInfo.Vdpa.Path)
							}
							if iface.NetworkStatus.DeviceInfo.Vdpa.PciAddress != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Vdpa.PciAddress =
									C.CString(iface.NetworkStatus.DeviceInfo.Vdpa.PciAddress)
							}
							if iface.NetworkStatus.DeviceInfo.Vdpa.PfPciAddress != "" {
								c_ifaceResp_pIface[j].NetworkStatus.DeviceInfo.Vdpa.PfPciAddress =
									C.CString(iface.NetworkStatus.DeviceInfo.Vdpa.PfPciAddress)
							}
						} else {
							_ = logging.Errorf("Error: type set to vDPA, but no associated DeviceInfo data")
						}
					default:
						c_ifaceResp_pIface[j].DeviceType = C.NETUTIL_TYPE_UNKNOWN
					}
				}

				c_ifaceRsp.numIfacePopulated++

				j++
			} else {

				_ = logging.Errorf("InterfaceResponse struct not sized properly."+
					"At Interface %d.", i)

				return NETUTIL_ERRNO_SIZE_ERROR
			}
		}
		return NETUTIL_ERRNO_SUCCESS
	}
	_ = logging.Errorf("netlib.GetInterfaces() err: %+v", err)
	return NETUTIL_ERRNO_FAIL
}

//export SetConfig
func SetConfig(c_config *C.struct_AppNetutilConfig) int64 {
	err := setLogSettings(&c_config.log)

	return err
}

func setLogSettings(c_log *C.struct_LoggingData) int64 {

	if c_log.level < C.LOG_LEVEL_MAX {
		logging.SetLogLevel(logging.Level.String(logging.Level(c_log.level)))
	}

	if c_log.stderr <= 1 {
		stderr := true
		if c_log.stderr == 0 {
			stderr = false
		}
		logging.SetLogStderr(stderr)
	}

	if c_log.filename != nil {
		logging.SetLogFile(C.GoString(c_log.filename))
	}

	return NETUTIL_ERRNO_SUCCESS
}


func main() {}
