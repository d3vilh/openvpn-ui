#!/bin/bash
# VERSION 1.3 by d3vilh@github.com aka Mr. Philipp. Thanks bugsyb@github.com for all the efforts ;)
set -e  # Exit immediately if a command exits with a non-zero status. Set -x option for debugging

# Define the machine architecture
# PLATFORM="linux/amd64" # arm64v8 = "linux/arm64/v8", arm32v5 - "linux/arm/v5", arm32v7 - "linux/arm/v7", amd64 - "linux/amd64"
ARCH=$(uname -m)
case $ARCH in
  armv6*)
    PLATFORM="linux/arm/v5"
    #UIIMAGE="FROM arm32v5/debian:stable-slim"
    UIIMAGE="FROM arm32v6/alpine" #moving to unstable because it has easy-rsa v3.1.6 which supports cert renewal
    BEEIMAGE="FROM arm32v5/golang:1.21-bookworm"
    ;;
  armv7*)
    PLATFORM="linux/arm/v7"
    #UIIMAGE="FROM arm32v7/debian:stable-slim"
    UIIMAGE="FROM arm32v7/alpine" #moving to unstable because it has easy-rsa v3.1.6 which supports cert renewal
    BEEIMAGE="FROM arm32v7/golang:1.22.3-bookworm"
    ;;
  arm64*)
    PLATFORM="linux/arm64/v8"
    #UIIMAGE="FROM arm64v8/debian:stable-slim"
    UIIMAGE="FROM arm64v8/alpine" #moving to unstable because it has easy-rsa v3.1.6 which supports cert renewal
    BEEIMAGE="FROM golang:1.22.3-bookworm"
    ;;
  *)
    PLATFORM="linux/amd64"
    #UIIMAGE="FROM debian:stable-slim"
    UIIMAGE="FROM alpine" #moving to unstable because it has easy-rsa v3.1.6 which supports cert renewal
    BEEIMAGE="FROM golang:1.22.3-bookworm"
    ;;
esac

# Benchmarking the start time get
start_time=$(date +%s)

printf "\033[1;34mBuilding for\033[0m $ARCH ($PLATFORM) with: \n  \033[1;34mUI Image:\033[0m $UIIMAGE \n  \033[1;34mBeeGo Image:\033[0m $BEEIMAGE \n"
# Update Dockerfile based on platform
sed "s#FROM DEFINE-YOUR-ARCH#$UIIMAGE#g" Dockerfile > Dockerfile.tmp && mv Dockerfile.tmp Dockerfile
# Update Dockerfile-beego based on platform
sed "s#FROM DEFINE-YOUR-ARCH#$BEEIMAGE#g" Dockerfile-beego > Dockerfile-beego.tmp && mv Dockerfile-beego.tmp Dockerfile-beego
printf "Dockerfiles updated \n\033[1;34mBuilding Golang and Bee environment.\033[0m\n"

# Build golang & bee environment
docker build --platform=$PLATFORM -f Dockerfile-beego -t local/beego-v8 -t local/beego-v8:latest .
printf "\033[1;34mBuilding OpenVPN-UI and qrencode binaries.\033[0m\n"

# Run a beego-v8 container to build qrencode and execute bee pack
time docker run \
    -v "$PWD/../":/go/src/github.com/d3vilh/openvpn-ui \
    -e GO111MODULE='auto' \
    -e CGO_ENABLED=1 \
    --rm \
    -w /usr/src/myapp \
    local/beego-v8 \
sh -c "cd /go/src/github.com/d3vilh/openvpn-ui/ && go env -w GOFLAGS="-buildvcs=false" && bee version && CGO_ENABLED=1 CC=musl-gcc bee pack -exr='^vendor|^ace.tar.bz2|^data.db|^build|^README.md|^docs' && cd /app/qrencode && go build -o qrencode main.go && chmod +x /app/qrencode/qrencode && cp -p /app/qrencode/qrencode /go/src/github.com/d3vilh/openvpn-ui/"
printf "OpenVPN-UI and qrencode were built \n\033[1;34mBuilding OpenVPN-UI image.\033[0m\n"

# Build OpenVPN-UI image
QRFILE="qrencode"
UIFILE="openvpn-ui.tar.gz"
cp -f ../$QRFILE ./
cp -f ../$UIFILE ./

# Build openvpn-ui image
docker build -t local/openvpn-ui .
rm -f $UIFILE; rm -f $(basename $UIFILE); #rm -f $QRFILE; 
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
