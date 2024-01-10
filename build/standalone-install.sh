#!/bin/bash
# VERSION 0.1 by d3vilh@github.com aka Mr. Philipp.
#
# DRAFT! DO NOT USE IT!
#

# All the variables
GOVERSION="1.21.5"

# Check if Go is installed and the version is 1.21
go_version=$(go version 2>/dev/null | awk '{print $3}' | tr -d "go")
if [[ -z "$go_version" || "$go_version" < $GOVERSION ]]
then
    echo "Golang version 1.21 is not installed."
    read -p "Would you like to install it? (y/n) " -n 1 -r
    echo    # move to a new line
    if [[ $REPLY =~ ^[Yy]$ ]]
    then
        # Install Go
        wget https://golang.org/dl/go1.21.5.linux-amd64.tar.gz # x86_64 only, at the monment.
        sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
        export PATH=$PATH:/usr/local/go/bin
        echo "export PATH=$PATH:$(go env GOPATH)/bin" >> ~/.bashrc
        source ~/.bashrc
    else
        read -p "Would you like to continue without Golang 1.21 installation? (y/n) " -n 1 -r
        echo    # move to a new line
        if [[ ! $REPLY =~ ^[Yy]$ ]]
        then
            echo "Installation terminated by user."
            exit 1
        fi
    fi
fi

# Description
echo "This script will install OpenVPN-UI and all the dependencies on your local environment. No containers will be used."
echo "THIS SCRIPT IS IN DEVELOPMENT AND NOT READY FOR ANY USE."
# Ask for confirmation
read -p "Do you want to continue? (y/n)" -n 1 -r
echo    # move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    exit 1
fi

# Update your system
echo "Updating current enviroment with apt-get update"
sudo apt-get update -y

# Install necessary tools
echo "Installing dependencies (sed, gcc, git, musl-tools, easy-rsa, curl, jq, oathtool)"
sudo apt-get install -y sed gcc git musl-tools easy-rsa curl jq oathtool

echo "Downloading all go modules (go mod download)"
go mod download

echo "Installing BeeGo v2"
go install github.com/beego/bee/v2@develop
source ~/.bashrc # reload bashrc to get bee command
echo "Clonning qrencode into build directory"
git clone https://github.com/d3vilh/qrencode

# Set environment variables
export GO111MODULE='auto'
export CGO_ENABLED=1
export CC=musl-gcc 

# Change project directory
cd ../
echo "Building and packing OpenVPN-UI"
# Execute bee pack
go env -w GOFLAGS="-buildvcs=false"
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

printf "\033[1;34mAll done.\033[0m\n"