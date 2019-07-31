package apputil

import (
	"github.com/containernetworking/cni/pkg/types"
)

const (
	INTERFACE_TYPE_PCI = "pci"
	INTERFACE_TYPE_VHOST = "vhost"
)

type CPUResponse struct {
	CPUSet	string	`json:"cpuset,omitempty"`
}

type EnvResponse struct {
	Envs	map[string]string
}

type ResourceResponse struct {
}

type NetworkInterfaceResponse struct {
	Interface	[]*NetworkInterface
}

type NetworkInterface struct {
	Type	string		`json:"type,omitempty"`
	ID	string		`json:"id,omitempty"`
}

type NetworkStatusResponse struct{
	Status	[]*NetworkStatus	`json:"status,omitempty"`
}

type NetworkStatus struct {
	Name		string		`json:"name,omitempty"`
	Interface	string		`json:"interface,omitempty"`
	IPs		[]string	`json:"ips,omitempty"`
	Mac		string		`json:"mac,omitempty"`
}

type MultusNetworkStatus struct {
	Name      string    `json:"name"`
	Interface string    `json:"interface,omitempty"`
	IPs       []string  `json:"ips,omitempty"`
	Mac       string    `json:"mac,omitempty"`
	Default   bool      `json:"default,omitempty"`
	DNS       types.DNS `json:"dns,omitempty"`
}

