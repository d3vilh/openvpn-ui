#!/bin/bash -e

#CA_NAME=LocalCA
#SERVER_NAME=server
#EASY_RSA=/usr/share/easy-rsa

#mkdir -p /etc/openvpn/keys
#touch /etc/openvpn/keys/index.txt
#echo 01 > /etc/openvpn/keys/serial
cp -f /opt/scripts/vars.template /etc/openvpn/pki/vars

#$EASY_RSA/clean-all
#source /etc/openvpn/keys/vars
#export KEY_NAME=$CA_NAME
#echo "Generating CA cert"
#$EASY_RSA/build-ca
#export EASY_RSA="${EASY_RSA:-.}"

#$EASY_RSA/pkitool --initca $*

#export KEY_NAME=$SERVER_NAME

#echo "Generating server cert"
#$EASY_RSA/build-key-server $SERVER_NAME
#$EASY_RSA/pkitool --server $SERVER_NAME
