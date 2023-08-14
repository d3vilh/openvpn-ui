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

# Benchmarking the start time get
start_time=$(date +%s)

printf "\033[1;34mBuilding for\033[0m $ARCH ($PLATFORM) with: \n  \033[1;34mUI Image:\033[0m $UIIMAGE \n  \033[1;34mBeeGo Image:\033[0m $BEEIMAGE \n"
# Update Dockerfile based on platform
sed -i "s#FROM DEFINE-YOUR-ARCH#$UIIMAGE#g" Dockerfile
# Update Dockerfile-beego based on platform
sed -i "s#FROM DEFINE-YOUR-ARCH#$BEEIMAGE#g" Dockerfile-beego
printf "Dockerfiles updated \n\033[1;34mBuilding Golang and Bee enviroment.\033[0m\n"

# Build golang & bee environment
docker build --platform=$PLATFORM -f Dockerfile-beego -t local/beego-v8 -t local/beego-v8:latest .
printf "\033[1;34mBuilding OpenVPN-UI binary.\033[0m\n"
./openvpn-ui-pack2.sh

printf "OpenVPN-UI built \n\033[1;34mBuilding OpenVPN-UI image.\033[0m\n"
# Build OpenVPN-UI image
PKGFILE="openvpn-ui.tar.gz"
cp -f ../$PKGFILE ./

docker build -t local/openvpn-ui .
rm -f $PKGFILE; rm -f $(basename $PKGFILE)
printf "\033[1;34mAll done.\033[0m\n"
# Benchmarking the end time record
end_time=$(date +%s)

# Calculate the execution time in seconds
execution_time=$((end_time - start_time))

# Calculate the execution time in minutes and seconds
minutes=$((execution_time / 60))
seconds=$((execution_time % 60))

# Print the execution time in mm:ss format
printf "\033[1;34mExecution time: %02d:%02d\033[0m (%d sec)\n" $minutes $seconds $execution_time