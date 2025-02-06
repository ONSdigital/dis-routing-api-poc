#!/bin/bash -eux

pushd dis-routing-api-poc
  make test-component
popd
