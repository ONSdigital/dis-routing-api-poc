#!/bin/bash -eux

go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5

pushd dis-routing-api-poc
  make lint
popd
