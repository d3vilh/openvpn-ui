#!/bin/bash -x
# Exit immediately if a command exits with a non-zero status
set -e

# Run a docker container with the specified volume and environment variable, and execute bee pack
time docker run \
    -v "$PWD/../":/go/src/github.com/d3vilh/openvpn-ui \
    -e GO111MODULE='auto' \
    --rm \
    -w /usr/src/myapp \
    local/beego-v8 \
    sh -c "cd /go/src/github.com/d3vilh/openvpn-ui/ && go env -w GOFLAGS="-buildvcs=false" && bee version && bee pack -exr='^vendor|^data.db|^build|^README.md|^docs'"