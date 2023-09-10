#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e
CERT_NAME=$1
CERT_IP=$2
CERT_SERIAL=$3
if [ -n "$1" ]; then
  export EASYRSA_BATCH=1
  # Temp WA for reneval
  #mv /usr/share/easy-rsa/pki/vars /usr/share/easy-rsa/pki/vars.bak
  # Renew certificate.
  echo "Renew certificate: $CERT_NAME with localip: $CERT_IP and serial: $CERT_SERIAL"
  cd /usr/share/easy-rsa
  ./easyrsa renew "$CERT_NAME" nopass  #as of now only nopass is supported
  
  # Fix for new /name in index.txt (adding name and ip to the last line)
  sed -i'.bak' "$ s/$/\/name=${CERT_NAME}\/LocalIP=${CERT_IP}/" /usr/share/easy-rsa/pki/index.txt
  
  # Fix for new /name in index.txt (adding name and ip to the last line)
  #sed -i'.bak' "s/\/name=${1}\/.*//" /usr/share/easy-rsa/pki/index.txt
  #./easyrsa revoke-renewed "$1"
  # Fix for new /name in index.txt (adding name and ip to the last line)
 
  #chmod +r ./pki/crl.pem
  #mv /usr/share/easy-rsa/pki/vars.back /usr/share/easy-rsa/pki/vars
  echo 'All Done, Sudar!'
  echo 
else
  echo "Invalid input argument: $CERT_NAME"
  exit 1
fi