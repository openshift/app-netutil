package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	netlib "github.com/openshift/app-netutil/lib/v1alpha"
)

func main() {
	flag.Parse()
	glog.Infof("starting sample application")

	for {
		cpuResponse, err := netlib.GetCPUInfo()
		if err != nil {
			glog.Errorf("Error calling netlib.GetCPUInfo: %v", err)
			return
		}
		glog.Infof("netlib.GetCPUInfo Response: %s", cpuResponse.CPUSet)

		envResponse, err := netlib.GetEnv()
		if err != nil {
			glog.Errorf("Error calling netlib.GetEnv: %v", err)
			return
		}
		glog.Infof("netlib.GetEnv Response:")
		for key, value := range envResponse.Envs {
			fmt.Printf("| %-25s|: %s\n", key, value)
		}

		netResponse, err := netlib.GetNetworkStatus()
		if err != nil {
			glog.Errorf("Error calling netlib.GetNetworkStatus: %v", err)
			return
		}
		glog.Infof("netlib.GetNetworkStatus Response:")
		for index, s := range netResponse.Status {
			fmt.Printf("| %-25d|: %+v\n", index, s)
		}

		intResponse, err := netlib.GetNetworkInterface("pci")
		if err != nil {
			glog.Errorf("Error calling netlib.GetInterfaces: %v", err)
			return
		}
		glog.Infof("netlib.GetInterfaces Response:")
		for index, iface := range intResponse.Interface {
			fmt.Printf("| %-25d|: %+v: %+v\n", index, iface.Type, iface.ID)
		}
		time.Sleep(1 * time.Minute)
	}
	return
}
