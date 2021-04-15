package apputil

import (
	nritypes "github.com/intel/network-resources-injector/pkg/types"
)

var downwardAPIMountPath = nritypes.DownwardAPIMountPath

func GetDownwardAPIMountPath() string {
	return downwardAPIMountPath
}

func SetDownwardAPIMountPath(mntPnt string) {
	downwardAPIMountPath = mntPnt
}
