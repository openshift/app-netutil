package apputil

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/golang/glog"
)

const (
	netStatusPath = "/etc/podinfo/annotations"
	netStatusKey = "k8s.v1.cni.cncf.io/networks-status"
)

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
