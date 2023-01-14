#!/bin/bash

set -e

PKGFILE=openvpn-ui.tar.gz

cp -f ../$PKGFILE ./

docker build -t d3vilh/openvpn-ui .

rm -f $PKGFILE
