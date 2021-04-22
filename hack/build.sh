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

GIT_COMMIT=$(git log -1 --pretty=format:%h)
GO_COMMIT_VAR="github.com/openshift/app-netutil/lib/v1alpha.GitCommit"

GIT_VERSION=$(git tag | sort -V | tail -1)
GO_VERSION_VAR="github.com/openshift/app-netutil/lib/v1alpha.AppNetutilVersion"

go install \
  -ldflags "-X '${GO_COMMIT_VAR}=${GIT_COMMIT}' -X '${GO_VERSION_VAR}=${GIT_VERSION}'" \
  "$@" \
  ${REPO_PATH}/samples/go_app
