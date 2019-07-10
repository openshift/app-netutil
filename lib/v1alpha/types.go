package apputil

import (
	"github.com/containernetworking/cni/pkg/types"
)

type CPUResponse struct {
	CPUSet	string	`json:"cpuset,omitempty"`
}

type EnvResponse struct {
	Envs	map[string]string
}

type ResourceResponse struct {
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

