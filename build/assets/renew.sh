#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

if [ "$1" = "renew" ]; then
  # Renew certificate. works 30days prior expiration.
  echo 'Renew certificate...'
  cd /usr/share/easy-rsa
  ./easyrsa renew "$2"
  chmod +r ./pki/crl.pem
  echo 'Done!'
  echo 
else
  echo "Invalid input argument: $1"
  exit 1
fi