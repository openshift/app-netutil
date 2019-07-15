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
		}
		glog.Infof("netlib.GetCPUInfo Response: %s", cpuResponse.CPUSet)

		envResponse, err := netlib.GetEnv()
		if err != nil {
			glog.Errorf("Error calling netlib.GetEnv: %v", err)
		}
		glog.Infof("netlib.GetEnv Response:")
		for key, value := range envResponse.Envs {
			fmt.Printf("| %-25s|: %s\n", key, value)
		}

		netResponse, err := netlib.GetNetworkStatus()
		if err != nil {
			glog.Errorf("Error calling netlib.GetNetworkStatus: %v", err)
		}
		glog.Infof("netlib.GetNetworkStatus Response:")
		for index, s := range netResponse.Status {
			fmt.Printf("| %-25d|: %+v\n", index, s)
		}
		time.Sleep(1 * time.Minute)
	}
	return
}
