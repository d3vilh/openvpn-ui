#!/bin/bash
# VERSION 0.3 by d3vilh@github.com aka Mr. Philipp.
#

# All the variables
GOVERSION="1.22.3"

# Description
echo "This script will install OpenVPN-UI and all the dependencies on your local environment. No containers will be used."
# Ask for confirmation
read -p "Do you want to continue? (y/n)" -n 1 -r
echo    # move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    exit 1
fi

# Check if Go is installed and the version is supported
go_version=$(go version 2>/dev/null | awk '{print $3}' | tr -d "go")
if [[ -z "$go_version" || "$go_version" < $GOVERSION ]]
then
    echo "Golang version ${GOVERSION} is not installed."
    read -p "Would you like to install it? (y/n) " -n 1 -r
    echo    # move to a new line
    if [[ $REPLY =~ ^[Yy]$ ]]
    then
        # Check the architecture of the machine
        arch=$(uname -m)
        if [[ "$arch" == "x86_64" ]]; then
            # Install Go for x86_64
            wget https://golang.org/dl/go${GOVERSION}.linux-arm64.tar.gz
            sudo tar -C /usr/local -xzf go${GOVERSION}.linux-arm64.tar.gz
        elif [[ "$arch" == "aarch64" ]]; then
            # Install Go for arm64
            wget https://golang.org/dl/go${GOVERSION}.linux-arm64.tar.gz
            sudo tar -C /usr/local -xzf go${GOVERSION}.linux-arm64.tar.gz
        elif [[ "$arch" == "armv7l" ]]; then
            # Install Go for armv7l
            wget https://golang.org/dl/go${GOVERSION}.linux-armv6l.tar.gz
            sudo tar -C /usr/local -xzf go${GOVERSION}.linux-armv6l.tar.gz
        elif [[ "$arch" == "armv6l" ]]; then
            # Install Go for armv6l
            wget https://golang.org/dl/go${GOVERSION}.linux-armv6l.tar.gz
            sudo tar -C /usr/local -xzf go${GOVERSION}.linux-armv6l.tar.gz
        else
            echo "Unsupported architecture."
            exit 1
        fi
        export PATH=$PATH:/usr/local/go/bin
        echo "export PATH=$PATH:$(go env GOPATH)/bin" >> ~/.bashrc
        source ~/.bashrc
    else
        read -p "Would you like to continue without Golang ${GOVERSION} installation? (y/n) " -n 1 -r
        echo    # move to a new line
        if [[ ! $REPLY =~ ^[Yy]$ ]]
        then
            echo "Installation terminated by user."
            exit 1
        fi
    fi
fi

# Update your system
read -p "Would you like to run apt-get update? (y/n) " -n 1 -r
echo    # move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
    # Update your system
    echo "Updating current environment with apt-get update"
    sudo apt-get update -y
fi

# Install necessary tools
read -p "Would you like to install the dependencies (sed, gcc, git, musl-tools, easy-rsa, curl, jq, oathtool)? (y/n) " -n 1 -r
echo    # move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
    # Install necessary tools
    echo "Installing dependencies (sed, gcc, git, musl-tools, easy-rsa, curl, jq, oathtool)"
    sudo apt-get install -y sed gcc git musl-tools easy-rsa curl jq oathtool
fi

# Go Modules download
read -p "Would you like to download all necessary Go modules? (y/n) " -n 1 -r
echo    # move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
    # Download all Go modules
    echo "Downloading all Go modules (go mod download)"
    go mod download
fi

read -p "Would you like to install Beego v2? (y/n) " -n 1 -r
echo    # move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
    # Install Beego
    echo "Installing BeeGo v2"
    go install github.com/beego/bee/v2@develop
fi

# Installing OpenVPN-UI and qrencode
read -p "Would you like to build OpenVPN-UI and qrencode binaries? (y/n) " -n 1 -r
echo    # move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
    # Install OpenVPN-UI and qrencode
    echo "Installing OpenVPN-UI and qrencode"
    source ~/.bashrc # reload bashrc to get bee command
    echo "Cloning qrencode into build directory"
    git clone https://github.com/d3vilh/qrencode

    # Set environment variables
    export GO111MODULE='auto'
    export CGO_ENABLED=1
    export CC=musl-gcc 

    # Packing openvpn-ui
    cd ../
    echo "Building and packing OpenVPN-UI"
    # Execute bee pack
    export PATH=$PATH:$(go env GOPATH)/bin
    go env -w GOFLAGS="-buildvcs=false"
    source ~/.bashrc
    bee version
    bee pack -exr='^vendor|^ace.tar.bz2|^data.db|^build|^README.md|^docs'

    # Build qrencode
    echo "Building qrencode"
    cd build/qrencode
    go build -o qrencode main.go
    chmod +x qrencode
    echo "Moving qrencode to GOPATH"
    mv qrencode $(go env GOPATH)/bin
    cd ../
fi

printf "\033[1;34mAll done.\033[0m\n"