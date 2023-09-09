#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

NAME=$1
SERIAL=$2

# .ovpn file path
DEST_FILE_PATH="/etc/openvpn/clients/$NAME.ovpn"

# Check if .ovpn file exists
if [[ ! -f $DEST_FILE_PATH ]]; then
    echo "User not found."
    exit 1
fi

# Fix index.txt by removing everything after pattern "/name=$1" in the line
# Fix for https://github.com/d3vilh/openvpn-ui/issues/5 by shuricksumy@github
#sed -i'.bak' "s/\/name=${name}\/.*//" /usr/share/easy-rsa/pki/index.txt

# removing the end of the line starting from /name=$NAME for the line that matches the $serial pattern
sed  -i'.bak' "/$SERIAL/s/\/name=$NAME.*//" /usr/share/easy-rsa/pki/index.txt

export EASYRSA_BATCH=1 # see https://superuser.com/questions/1331293/easy-rsa-v3-execute-build-ca-and-gen-req-silently

echo 'Revoke certificate...'

cd /usr/share/easy-rsa
# Revoke certificate
./easyrsa revoke "$NAME"

echo 'Create new Create certificate revocation list (CRL)...'
./easyrsa gen-crl
chmod +r ./pki/crl.pem

echo 'Done!'
echo 'If you want to disconnect the user please restart the service using docker-compose restart openvpn.'
