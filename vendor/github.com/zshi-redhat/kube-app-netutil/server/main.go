package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	api "github.com/zshi-redhat/kube-app-netutil/apis/v1alpha"
)

const (
	endpoint = "/etc/netutil/net.sock"
	cpusetPath = "/sys/fs/cgroup/cpuset/cpuset.cpus"
	netStatusPath = "/etc/podinfo/annotations"
	netStatusKey = "k8s.v1.cni.cncf.io/networks-status"
)

type NetUtilServer struct {
	grpcServer	*grpc.Server
}

type MultusNetworkStatus struct {
	Name      string    `json:"name"`
	Interface string    `json:"interface,omitempty"`
	IPs       []string  `json:"ips,omitempty"`
	Mac       string    `json:"mac,omitempty"`
	Default   bool      `json:"default,omitempty"`
	DNS       types.DNS `json:"dns,omitempty"`
}

func newNetutilServer() *NetUtilServer {
	return &NetUtilServer{
		grpcServer: grpc.NewServer(),
	}
}

func (ns *NetUtilServer) start() error {
	glog.Infof("starting netutil server at: %s\n", endpoint)
	lis, err := net.Listen("unix", endpoint)
	if err != nil {
		glog.Errorf("Error creating netutil gRPC service: %v", err)
		return err
	}

	api.RegisterNetUtilServer(ns.grpcServer, ns)
	go ns.grpcServer.Serve(lis)

	// Wait for server to start
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		glog.Errorf("Error starting netutil server: %v", err)
		return err
	}
	glog.Infof("netutil server start serving")
	conn.Close()
	return nil
}

func (ns *NetUtilServer) stop() error {
	glog.Infof("stopping netutil server")
	if ns.grpcServer != nil {
		ns.grpcServer.Stop()
		ns.grpcServer = nil
	}
	err := os.Remove(endpoint)
	if err != nil && !os.IsNotExist(err) {
		glog.Errorf("Error cleaning up socket file")
	}
	return nil
}

func (ns *NetUtilServer) GetCPUInfo(ctx context.Context, rqt *api.CPURequest) (*api.CPUResponse, error) {
	path := filepath.Join("/proc", rqt.Pid, "root", cpusetPath)
	glog.Infof("getting cpuset from path: %s", path)
	cpus, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Errorf("Error getting cpuset info: %v", err)
		return nil, err
	}
	return &api.CPUResponse{Cpuset: string(bytes.TrimSpace(cpus))}, nil
}

func (ns *NetUtilServer) GetEnv(ctx context.Context, rqt *api.EnvRequest) (*api.EnvResponse, error) {
	path := filepath.Join("/proc", rqt.Pid, "environ")
	glog.Infof("getting environment variables from path: %s", path)
	file, err := os.Open(path)
	if err != nil {
		glog.Errorf("Error openning proc environ file: %v", err)
		return nil, err
	}
	defer file.Close()

	envAttrs := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		envs := strings.Split(string(line), "\x00")
		for _, e := range envs {
			parts := strings.Split(string(e), "=")
			if len(parts) == 2 {
				envAttrs[parts[0]] = parts[1]
			}
		}
	}
	return &api.EnvResponse{Envs: envAttrs}, nil
}

func (ns *NetUtilServer) GetNetworkStatus(ctx context.Context, rqt *api.NetworkStatusRequest) (*api.NetworkStatusResponse, error) {
	glog.Infof("getting network status from path: %s", netStatusPath)
	file, err := os.Open(netStatusPath)
	if err != nil {
		glog.Errorf("Error openning network status file: %v", err)
		return nil, err
	}
	defer file.Close()

	data := []MultusNetworkStatus{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		status := strings.Split(string(line), "\n")
		for _, s := range status {
			parts := strings.Split(string(s), "=")
			if len(parts) == 2 {
				if parts[0] == netStatusKey {
					parts[1] = strings.Replace(string(parts[1]), "\\n", "", -1)
					parts[1] = strings.Replace(string(parts[1]), "\\", "", -1)
					parts[1] = strings.Replace(string(parts[1]), " ", "", -1)
					parts[1] = string(parts[1][1:len(parts[1])-1])
					if err = json.Unmarshal([]byte(parts[1]), &data); err != nil {
						glog.Errorf("Error unmarshal multus network status: %v", err)
					}
				}
			}
		}
	}

	response := &api.NetworkStatusResponse{}
	for _, status := range data {
		response.Status = append(response.Status, &api.NetworkStatus{
				Name: status.Name,
				Interface: status.Interface,
				Mac: status.Mac,
				Ips: status.IPs,
			})
	}
	return response, nil
}

func (ns *NetUtilServer) GetResource(ctx context.Context, rqt *api.ResourceRequest) (*api.ResourceResponse, error) {
	return &api.ResourceResponse{}, nil
}

func main() {
	flag.Parse()
	ns := newNetutilServer()
	if ns == nil {
		glog.Errorf("Error initializing netutil manager")
		return
	}
	err := ns.start()
	if err != nil {
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-sigCh:
		glog.Infof("signal received, shutting down", sig)
		ns.stop()
		return
	}
}
