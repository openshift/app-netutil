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
		glog.Infof("CALL netlib.GetCPUInfo:")
		cpuResponse, err := netlib.GetCPUInfo()
		if err != nil {
			glog.Errorf("Error calling netlib.GetCPUInfo: %v", err)
			return
		}
		glog.Infof("netlib.GetCPUInfo Response:")
		fmt.Printf("| CPU     |: %+v\n", cpuResponse.CPUSet)

		glog.Infof("CALL netlib.GetInterfaces:")
		ifaceResponse, err := netlib.GetInterfaces()
		if err != nil {
			glog.Errorf("Error calling netlib.GetInterfaces: %v", err)
			return
		}
		glog.Infof("netlib.GetInterfaces Response:")
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
