#!/bin/bash
#VERSION 1.1
# Exit immediately if a command exits with a non-zero status
set -e

# .ovpn file path
DEST_FILE_PATH="/etc/openvpn/clients/$1.ovpn"
INDEX_PATH="/usr/share/easy-rsa/pki/index.txt"

# Check if .ovpn file exists
if [[ ! -f $DEST_FILE_PATH ]]; then
    echo "User not found."
    exit 1
fi

# Define key serial number by keyname
# Fix for https://github.com/d3vilh/openvpn-ui/issues/5 by shuricksumy@github
STATUS_CH=$(grep -e ${1}$ -e${1}/ ${INDEX_PATH} | awk '{print $1}' | tr -d '\n')
if [[ $STATUS_CH = "V" ]]; then
    echo "Cert is VALID"
    CERT_SERIAL=$(grep ${1}/ ${INDEX_PATH} | awk '{print $3}' | tr -d '\n')
    echo "Will remove: ${CERT_SERIAL}"
else
    echo "Cert is REVOKED"
    CERT_SERIAL=$(grep ${1}$ ${INDEX_PATH} | awk '{print $4}' | tr -d '\n')
    echo "Will remove: ${CERT_SERIAL}"
fi

# Remove user from OpenVPN
rm -f /etc/openvpn/pki/certs_by_serial/$CERT_SERIAL.pem
rm -f /etc/openvpn/pki/issued/$1.crt
rm -f /etc/openvpn/pki/private/$1.key
rm -f /etc/openvpn/pki/reqs/$1.req
rm -f /etc/openvpn/clients/$1.ovpn

# Fix index.txt by removing the user from the list following the serial number
sed -i'.bak' "/${CERT_SERIAL}/d" $INDEX_PATH

echo 'Remove done!'
echo 'If you want to disconnect the user please restart the OpenVPN service or container.'
