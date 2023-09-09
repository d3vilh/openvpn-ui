#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

if [ -n "$1" ]; then
  export EASYRSA_BATCH=1
  # Temp WA for reneval
  #mv /usr/share/easy-rsa/pki/vars /usr/share/easy-rsa/pki/vars.bak
  # Renew certificate.
  echo "Renew certificate: $1 with localip: $2 and serial: $3"
  cd /usr/share/easy-rsa
  ./easyrsa renew "$1" nopass  #as of now only nopass is supported
  
  # Fix for new /name in index.txt (adding name and ip to the last line)
  sed -i'.bak' "$ s/$/\/name=${1}\/LocalIP=${2}/" /usr/share/easy-rsa/pki/index.txt
  
  #sed -i'.bak' "s/\/name=${1}\/.*//" /usr/share/easy-rsa/pki/index.txt
  #./easyrsa revoke-renewed "$1"
  # Fix for new /name in index.txt (adding name and ip to the last line)
 
  chmod +r ./pki/crl.pem
  #mv /usr/share/easy-rsa/pki/vars.back /usr/share/easy-rsa/pki/vars
  echo 'All Done, Sudar!'
  echo 
else
  echo "Invalid input argument: $1"
  exit 1
fi