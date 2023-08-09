#!/bin/bash -x
# Exit immediately if a command exits with a non-zero status
set -e

PLATFORM="linux/arm64/v8" # amd64 - "linux/amd64", arm32v5 - "linux/arm/v5", arm32v7 - "linux/arm/v7"
# Set the current directory
#OD=$PWD

# go to script folder, required for below steps
#CURDIR=${OD}/$(dirname $0)
#cd "${CURDIR}"

# Build golang & bee environment
docker build --platform=$PLATFORM -f Dockerfile-beego -t local/beego-v8 -t local/beego-v8:latest .
./openvpn-ui-pack2.sh

# Build OpenVPN-UI image
PKGFILE="openvpn-ui.tar.gz"
cp -f ../$PKGFILE ./

docker build -t local/openvpn-ui .
rm -f $PKGFILE; rm -f $(basename $PKGFILE)
