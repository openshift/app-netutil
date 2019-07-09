package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	api "github.com/zshi-redhat/kube-app-netutil/apis/v1alpha"
)

const (
	endpoint = "/etc/netutil/net.sock"
)

func main() {
	flag.Parse()
	glog.Infof("starting netutil client")

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	defer conn.Close()

	c := api.NewNetUtilClient(conn)
	pid := os.Getpid()
	glog.Infof("os.Getpid() returns: %v", pid)
	cpuResponse, err := c.GetCPUInfo(context.Background(), &api.CPURequest{Pid: strconv.Itoa(pid)})
	if err != nil {
		glog.Errorf("Error calling GetCPUInfo: %v", err)
	}
	glog.Infof("GetCPUInfo Response from NetUtilServer: %s", cpuResponse.Cpuset)

	envResponse, err := c.GetEnv(context.Background(), &api.EnvRequest{Pid: strconv.Itoa(pid)})
	if err != nil {
		glog.Errorf("Error calling GetEnv: %v", err)
	}
	glog.Infof("GetEnv Response from NetUtilServer:")
	for key, value := range envResponse.Envs {
		fmt.Printf("| %-25s|: %s\n", key, value)
	}

	netResponse, err := c.GetNetworkStatus(context.Background(), &api.NetworkStatusRequest{Pid: strconv.Itoa(pid)})
	if err != nil {
		glog.Errorf("Error calling GetNetworkStatus: %v", err)
	}
	glog.Infof("GetEnv Response from NetUtilServer:")
	for index, s := range netResponse.Status {
		fmt.Printf("| %-25d|: %v\n", index, s)
	}
	return
}
