package apputil

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/golang/glog"

	"github.com/openshift/app-netutil/pkg/types"
)

//
// API Functions
//
func GetHugepages() (*types.HugepagesResponse, error) {
	var request int64
	var limit int64

	reqPath := filepath.Join(types.DownwardAPIMountPath, types.HugepagesRequestPath)
	glog.Infof("GetHugepages: Open %s", reqPath)
	requestStr, reqErr := ioutil.ReadFile(reqPath)
	if reqErr != nil {
		glog.Infof("Error getting %s info: %v", reqPath, reqErr)
		request = 0
	} else {
		request, reqErr = strconv.ParseInt(string(bytes.TrimSpace(requestStr)), 10, 64)
		if reqErr != nil {
			glog.Infof("Error converting limit \"%s\": %v", requestStr, reqErr)
		}
	}

	limPath := filepath.Join(types.DownwardAPIMountPath, types.HugepagesLimitPath)
	glog.Infof("GetHugepages: Open %s", limPath)
	limitStr, limErr := ioutil.ReadFile(limPath)
	if limErr != nil {
		glog.Infof("Error getting %s info: %v", limPath, limErr)
		limit = 0
	} else {
		limit, limErr = strconv.ParseInt(string(bytes.TrimSpace(limitStr)), 10, 64)
		if limErr != nil {
			glog.Infof("Error converting limit \"%s\": %v", limitStr, limErr)
		}
	}

	if reqErr != nil && limErr != nil {
		return nil, reqErr
	}
	return &types.HugepagesResponse{
		Request: request,
		Limit:   limit,
	}, nil
}
