#!/bin/bash
#VERSION 1.2 by d3vilh@github.com aka Mr. Philipp
# Exit immediately if a command exits with a non-zero status
set -e

#Variables
CERT_NAME=$1
CERT_IP=$2
CERT_SERIAL=$3
EASY_RSA=$(grep -E "^EasyRsaPath\s*=" ../openvpn-ui/conf/app.conf | cut -d= -f2 | tr -d '"' | tr -d '[:space:]')
echo 'EasyRSA path: $EASY_RSA'

# Check if 2FA was specified. If not - set to none.
if [ -z "$TFA_NAME" ]; then
    TFA_NAME="none"
fi

if [ -n "$1" ]; then
  export EASYRSA_BATCH=1

  # Renew certificate.
  echo "Renew certificate: $CERT_NAME with $TFA_NAME with localip: $CERT_IP and serial: $CERT_SERIAL"
  cd $EASY_RSA
  ./easyrsa renew "$CERT_NAME" nopass  #as of now only nopass is supported
  
  # Fix for new /name in index.txt (adding name and ip to the last line)
  sed -i'.bak' "$ s/$/\/name=${CERT_NAME}\/LocalIP=${CERT_IP}\/2FAName=${TFA_NAME}/" $EASY_RSA/pki/index.txt
  echo 'All Done, hlopche!'
else
  echo "Invalid input argument: $CERT_NAME Exiting."
  exit 1
fi