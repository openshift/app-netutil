set -e

ORG_PATH="github.com/zshi-redhat"
REPO_PATH="${ORG_PATH}/kube-app-netutil"

if [ ! -h gopath/src/${REPO_PATH} ]; then
        mkdir -p gopath/src/${ORG_PATH}
        ln -s ../../../.. gopath/src/${REPO_PATH} || exit 255 
fi

export GOBIN=${PWD}/bin
export GOPATH=${PWD}/gopath
export CGO_ENABLED=0

go install "$@" ${REPO_PATH}/server
go install "$@" ${REPO_PATH}/client
