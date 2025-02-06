#!/bin/bash -eux

pushd dis-routing-api-poc
  make build
  cp build/dis-routing-api-poc Dockerfile.concourse ../build
popd
