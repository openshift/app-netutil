package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	nettypes "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	netlib "github.com/openshift/app-netutil/lib/v1alpha"
)

func main() {
	flag.Parse()
	glog.Infof("starting sample application")

	for {
		// Test GetCPUInfo()
		glog.Infof("CALL netlib.GetCPUInfo:")
		cpuResponse, err := netlib.GetCPUInfo()
		if err != nil {
			glog.Errorf("Error calling netlib.GetCPUInfo: %v", err)
			return
		}
		fmt.Printf("|CPU      |: %+v\n", cpuResponse.CPUSet)

		// Test GetHugepages()
		glog.Infof("CALL netlib.GetHugepages:")
		hugepagesResponse, err := netlib.GetHugepages()
		if err != nil {
			glog.Infof("Error calling netlib.GetHugepages: %v", err)
		} else {
			fmt.Printf("|HugePage |: MyContainerName=%+v\n", hugepagesResponse.MyContainerName)
			if len(hugepagesResponse.Hugepages) != 0 {
				for index, hugepagesData := range hugepagesResponse.Hugepages {
					fmt.Printf("| %-8d|: ContainerName=%+v  Request 1G=%+v 2M=%+v Ukn=%+v  Limit 1G=%+v 2M=%+v Ukn=%+v\n",
						index,
						hugepagesData.ContainerName,
						hugepagesData.Request1G,
						hugepagesData.Request2M,
						hugepagesData.Request,
						hugepagesData.Limit1G,
						hugepagesData.Limit2M,
						hugepagesData.Limit)
				}
			}
		}

		// Test GetInterfaces()
		glog.Infof("CALL netlib.GetInterfaces:")
		ifaceResponse, err := netlib.GetInterfaces()
		if err != nil {
			glog.Errorf("Error calling netlib.GetInterfaces: %v", err)
			return
		}
		fmt.Printf("|Iface    |:\n")
		for index, iface := range ifaceResponse.Interface {
			fmt.Printf("| %-8d|: IfName=%+v  DeviceType=%+v  Network=%+v  Default=%+v\n",
				index, iface.NetworkStatus.Interface, iface.DeviceType,
				iface.NetworkStatus.Name, iface.NetworkStatus.Default)
			fmt.Printf("|         |:   MAC=%+v  IPs=%+v  DNS=%+v\n",
				iface.NetworkStatus.Mac, iface.NetworkStatus.IPs, iface.NetworkStatus.DNS)

			if iface.NetworkStatus.DeviceInfo != nil {
				switch iface.NetworkStatus.DeviceInfo.Type {
				case nettypes.DeviceInfoTypePCI:
					if iface.NetworkStatus.DeviceInfo.Pci != nil {
						fmt.Printf("|         |:   %+v  %+v  PCI=%+v  PF=%+v  Vhostnet=%+v  RdmaDevice=%+v\n",
							iface.NetworkStatus.DeviceInfo.Type,
							iface.NetworkStatus.DeviceInfo.Version,
							iface.NetworkStatus.DeviceInfo.Pci.PciAddress,
							iface.NetworkStatus.DeviceInfo.Pci.PfPciAddress,
							iface.NetworkStatus.DeviceInfo.Pci.Vhostnet,
							iface.NetworkStatus.DeviceInfo.Pci.RdmaDevice)
					}
				case nettypes.DeviceInfoTypeVHostUser:
					if iface.NetworkStatus.DeviceInfo.VhostUser != nil {
						fmt.Printf("|         |:   %+v  %+v  Mode=%+v  Path=%+v\n",
							iface.NetworkStatus.DeviceInfo.Type,
							iface.NetworkStatus.DeviceInfo.Version,
							iface.NetworkStatus.DeviceInfo.VhostUser.Mode,
							iface.NetworkStatus.DeviceInfo.VhostUser.Path)
					}
				case nettypes.DeviceInfoTypeMemif:
					if iface.NetworkStatus.DeviceInfo.Memif != nil {
						fmt.Printf("|         |:   %+v  %+v  Role=%+v  mode=%+v  Path=%+v\n",
							iface.NetworkStatus.DeviceInfo.Type,
							iface.NetworkStatus.DeviceInfo.Version,
							iface.NetworkStatus.DeviceInfo.Memif.Role,
							iface.NetworkStatus.DeviceInfo.Memif.Mode,
							iface.NetworkStatus.DeviceInfo.Memif.Path)
					}
				case nettypes.DeviceInfoTypeVDPA:
					if iface.NetworkStatus.DeviceInfo.Vdpa != nil {
						fmt.Printf("|         |:   %+v  %+v  ParentDevice=%+v  Driver=%+v  PciAddress=%+v  PfPciAddress=%+v  Path=%+v\n",
							iface.NetworkStatus.DeviceInfo.Type,
							iface.NetworkStatus.DeviceInfo.Version,
							iface.NetworkStatus.DeviceInfo.Vdpa.ParentDevice,
							iface.NetworkStatus.DeviceInfo.Vdpa.Driver,
							iface.NetworkStatus.DeviceInfo.Vdpa.PciAddress,
							iface.NetworkStatus.DeviceInfo.Vdpa.PfPciAddress,
							iface.NetworkStatus.DeviceInfo.Vdpa.Path)
					}
				}
			}
		}
		fmt.Printf("\n")

		time.Sleep(1 * time.Minute)
	}
}
