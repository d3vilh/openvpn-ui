#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

# Determine the machine architecture
# PLATFORM="linux/amd64" # arm64v8 = "linux/arm64/v8", arm32v5 - "linux/arm/v5", arm32v7 - "linux/arm/v7", amd64 - "linux/amd64"
ARCH=$(uname -m)
case $ARCH in
  armv6*)
    PLATFORM="linux/arm/v5"
    UIIMAGE="FROM arm32v5/debian:stable-slim"
    BEEIMAGE="FROM arm32v5/golang:1.20"
    ;;
  armv7*)
    PLATFORM="linux/arm/v7"
    UIIMAGE="FROM arm32v7/debian:stable-slim"
    BEEIMAGE="FROM arm32v7/golang:1.20"
    ;;
  aarch64*)
    PLATFORM="linux/arm64/v8"
    UIIMAGE="FROM arm64v8/debian:stable-slim"
    BEEIMAGE="FROM arm64v8/golang:1.20"
    ;;
  *)
    PLATFORM="linux/amd64"
    UIIMAGE="FROM debian:stable-slim"
    BEEIMAGE="FROM golang:1.20"
    ;;
esac

echo "Building for $ARCH ($PLATFORM) with UI Image $UIIMAGE and BeeGo Image $BEEIMAGE"
# Update Dockerfile based on platform
sed -i "s#FROM DEFINE-YOUR-ARCH#$UIIMAGE#g" Dockerfile
# Update Dockerfile-beego based on platform
sed -i "s#FROM DEFINE-YOUR-ARCH#$BEEIMAGE#g" Dockerfile-beego
echo "Dockerfiles updated \n Building Golang and Bee enviroment."

# Build golang & bee environment
docker build --platform=$PLATFORM -f Dockerfile-beego -t local/beego-v8 -t local/beego-v8:latest .
echo "Golang and Bee enviroment built \n Building OpenVPN-UI."
./openvpn-ui-pack2.sh

echo "OpenVPN-UI built \n Building OpenVPN-UI image."
# Build OpenVPN-UI image
PKGFILE="openvpn-ui.tar.gz"
cp -f ../$PKGFILE ./

docker build -t local/openvpn-ui .
rm -f $PKGFILE; rm -f $(basename $PKGFILE)
