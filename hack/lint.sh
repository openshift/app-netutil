#!/usr/bin/env bash

# SPDX-License-Identifier: Apache-2.0
# Copyright(c) 2021 Red Hat, Inc.

GO111MODULE=on ${GOPATH}/bin/golangci-lint run \
    --timeout=15m0s --verbose --print-resources-usage --modules-download-mode=vendor \
    && echo "lint OK!"
