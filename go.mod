module github.com/openshift/app-netutil

go 1.15

require (
	github.com/intel/network-resources-injector v0.0.0-20201215195952-4f073638930a
	github.com/intel/userspace-cni-network-plugin v0.0.0-20201116143459-807c52367c73
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v1.1.1-0.20201119153432-9d213757d22d
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace k8s.io/client-go => k8s.io/client-go v0.18.5
