#!/bin/bash
#VERSION 1.1
# Exit immediately if a command exits with a non-zero status
set -e

# .ovpn file path
DEST_FILE_PATH="/etc/openvpn/clients/$CERT_NAME.ovpn"
INDEX_PATH="/usr/share/easy-rsa/pki/index.txt"
CERT_SERIAL=$2
CERT_NAME=$1

# Check if .ovpn file exists
if [[ ! -f $DEST_FILE_PATH ]]; then
    echo "User not found."
    exit 1
fi

echo "Removing user: $CERT_NAME with Serial: $CERT_SERIAL"

# Define if cert is valid or revoked
STATUS_CH=$(grep -e ${CERT_NAME}$ -e${CERT_NAME}/ ${INDEX_PATH} | awk '{print $1}' | tr -d '\n')
if [[ $STATUS_CH = "V" ]]; then
    echo "Cert is VALID"
    echo "Will remove: ${CERT_SERIAL}"
else
    echo "Cert is REVOKED"
    echo "Will remove: ${CERT_SERIAL}"
fi

# Remove user from OpenVPN
rm -f /etc/openvpn/pki/certs_by_serial/$CERT_SERIAL.pem
rm -f /etc/openvpn/pki/issued/$CERT_NAME.crt
rm -f /etc/openvpn/pki/private/$CERT_NAME.key
rm -f /etc/openvpn/pki/reqs/$CERT_NAME.req
rm -f /etc/openvpn/clients/$CERT_NAME.ovpn

# Fix index.txt by removing the user from the list following the serial number
sed -i'.bak' "/${CERT_SERIAL}/d" $INDEX_PATH

echo 'Remove done!'
echo 'If you want to disconnect the user please restart the OpenVPN service or container.'
