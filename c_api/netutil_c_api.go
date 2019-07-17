package main
/*
#include <stdint.h>

// Mapped from app-netutil.lib/v1alpha/types.go

struct CPUResponse {
    char*    CPUSet;
};


struct EnvData {
    char*    Index;
    char*    Value;
};
// *pEnvs is an array of 'struct EnvData' that is allocated
// from the C program.
struct EnvResponse {
	int             netutil_num_envs;
	struct EnvData *pEnvs;
};


#define NETUTIL_NUM_IPS            10
#define NETUTIL_NUM_NETWORKSTATUS  10

struct NetworkStatus {
    char*    Name;
    char*    Interface;
    char*    IPs[NETUTIL_NUM_IPS];
    char*    Mac;
};
struct NetworkStatusResponse {
	struct NetworkStatus Status[NETUTIL_NUM_NETWORKSTATUS];
};
*/
import "C"
import "unsafe"

import (
	"flag"

	"github.com/golang/glog"

	netlib "github.com/openshift/app-netutil/lib/v1alpha"
)

const (
	cpusetPath = "/sys/fs/cgroup/cpuset/cpuset.cpus"
	netutil_num_ips = 10
	netutil_num_networkstatus = 10
)


//export GetCPUInfo
func GetCPUInfo(c_cpuResp *C.struct_CPUResponse) {
	flag.Parse()
	cpuRsp, err := netlib.GetCPUInfo()

	if err == nil {
		c_cpuResp.CPUSet = C.CString(cpuRsp.CPUSet)
	}
}

//export GetEnv
func GetEnv(c_envResp *C.struct_EnvResponse) {
	var j _Ctype_int

	flag.Parse()
	envRsp, err := netlib.GetEnv()

	if err == nil {
		j = 0

		// Map the input pointer to array of structures, c_envResp.pEnvs, to
		// a slice of the structures, c_envResp_pEnvs. Then the slice can be
		// indexed.
		c_envResp_pEnvs := (*[1 << 30]C.struct_EnvData)(unsafe.Pointer(c_envResp.pEnvs))[:c_envResp.netutil_num_envs:c_envResp.netutil_num_envs]
		for i, env := range envRsp.Envs {
			if j < c_envResp.netutil_num_envs {
				c_envResp_pEnvs[j].Index = C.CString(i)
				c_envResp_pEnvs[j].Value = C.CString(env)
				j++
			} else {
				glog.Errorf("EnvResponse struct not sized properly. At %d ENV Variables.", j)
				break
			}
		}
	}
}


//export GetNetworkStatus
func GetNetworkStatus(c_networkResp *C.struct_NetworkStatusResponse) {
	flag.Parse()
	networkStatusRsp, err := netlib.GetNetworkStatus()

	if err == nil {
		for i, networkStatus := range networkStatusRsp.Status {
			if i < netutil_num_networkstatus {
				c_networkResp.Status[i].Name = C.CString(networkStatus.Name)
				c_networkResp.Status[i].Interface = C.CString(networkStatus.Interface)
				c_networkResp.Status[i].Mac = C.CString(networkStatus.Mac)
				for j, ipaddr := range networkStatus.IPs {
					if j < netutil_num_ips {
						c_networkResp.Status[i].IPs[j] = C.CString(ipaddr)
					} else {
						glog.Errorf("NetworkStatusResponse IPs struct not sized properly. At %d IPs for Interface %d.", j, i)
						break
					}
				}
			} else {
				glog.Errorf("NetworkStatusResponse struct not sized properly. At %d Interfaces.", i)
				break
			}
		}
	}
}


func main() {}
