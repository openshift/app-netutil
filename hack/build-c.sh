set -e

ORG_PATH="github.com/openshift"
REPO_PATH="${ORG_PATH}/app-netutil"

if [ ! -h gopath/src/${REPO_PATH} ]; then
        mkdir -p gopath/src/${ORG_PATH}
        ln -s ../../../.. gopath/src/${REPO_PATH} || exit 255 
fi

export GOBIN=${PWD}/bin
export GOPATH=${PWD}/gopath
export CGO_ENABLED=1

GIT_COMMIT=$(git log -1 --pretty=format:%h)
GO_COMMIT_VAR="github.com/openshift/app-netutil/lib/v1alpha.GitCommit"

GIT_VERSION=$(git tag | sort -V | tail -1)
GO_VERSION_VAR="github.com/openshift/app-netutil/lib/v1alpha.AppNetutilVersion"

#go install "$@" ${REPO_PATH}/samples/go_app
go build \
  -ldflags "-X '${GO_COMMIT_VAR}=${GIT_COMMIT}' -X '${GO_VERSION_VAR}=${GIT_VERSION}'" \
  -o ${GOBIN}/libnetutil_api.so \
  -buildmode=c-shared \
  -v \
  ${REPO_PATH}/c_api

gcc \
  -I${GOBIN} \
  -L${GOBIN} \
  -Wall \
  -o ${GOBIN}/c_sample \
  ${GOPATH}/src/${REPO_PATH}/samples/c_app/app_sample.c -lnetutil_api
