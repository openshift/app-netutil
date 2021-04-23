# SPDX-License-Identifier: Apache-2.0
# Copyright(c) 2021 Red Hat, Inc.

set -e

ORG_PATH="github.com/openshift"
REPO_PATH="${ORG_PATH}/app-netutil"

if [ ! -h gopath/src/${REPO_PATH} ]; then
        mkdir -p gopath/src/${ORG_PATH}
        ln -s ../../../.. gopath/src/${REPO_PATH} || exit 255 
fi

export GOBIN=${PWD}/bin
export GOPATH=${PWD}/gopath
export CGO_ENABLED=0

GO_COMMIT_VAR="github.com/openshift/app-netutil/lib/v1alpha.GitCommit"
GIT_COMMIT=$(git log -1 --pretty=format:%h)
GIT_COMMIT_STR=" -X '${GO_COMMIT_VAR}=${GIT_COMMIT}'"

GIT_VERSION_STR=""
GO_VERSION_VAR="github.com/openshift/app-netutil/lib/v1alpha.AppNetutilVersion"
GIT_TAG=`git describe --tags`
if [ "$?" == 0 ]; then
   GIT_VERSION_STR=" -X '${GO_VERSION_VAR}=${GIT_TAG}'"
fi

go install \
  -ldflags "${GIT_COMMIT_STR}${GIT_VERSION_STR}" \
  "$@" \
  ${REPO_PATH}/samples/go_app
