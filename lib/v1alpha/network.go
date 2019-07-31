package apputil

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
)

const (
	netStatusPath = "/etc/podinfo/annotations"
	netStatusKey = "k8s.v1.cni.cncf.io/networks-status"
)

func GetNetworkInterface(intType string) (*NetworkInterfaceResponse, error) {
	glog.Infof("getting network interfaces")
	response := &NetworkInterfaceResponse{}

	switch intType {
	case INTERFACE_TYPE_PCI:
		envResponse, err := GetEnv()
		if err != nil {
			glog.Errorf("Error calling GetEnv from GetInterface: %v", err)
			return nil, err
		}
		for k, v := range envResponse.Envs {
			if strings.HasPrefix(k, "PCIDEVICE") {
				valueParts := strings.Split(string(v), ",")
				for _, id := range valueParts {
					response.Interface = append(response.Interface, &NetworkInterface{
						Type: k,
						ID: id,
					})
				}
			}
		}
	case INTERFACE_TYPE_VHOST:
		glog.Infof("Not implemented")
	case "":
		glog.Infof("Not implemented")
	default:
		return nil, fmt.Errorf("Unsupported interface type, values must be 'pci','vhost' or ''")
	}
	return response, nil
}

func GetNetworkStatus() (*NetworkStatusResponse, error) {
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

	response := &NetworkStatusResponse{}
	for _, status := range data {
		response.Status = append(response.Status, &NetworkStatus{
				Name: status.Name,
				Interface: status.Interface,
				Mac: status.Mac,
				IPs: status.IPs,
			})
	}
	return response, nil
}
