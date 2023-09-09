#!/bin/bash
#VERSION 1.1
# Exit immediately if a command exits with a non-zero status
set -e

# .ovpn file path
DEST_FILE_PATH="/etc/openvpn/clients/$CERT_NAME.ovpn"
INDEX="/usr/share/easy-rsa/pki/index.txt"
EASY_RSA="/usr/share/easy-rsa"
CERT_SERIAL=$2
CERT_NAME=$1

echo "Removing user: $CERT_NAME with Serial: $CERT_SERIAL"

# Define if cert is valid or revoked
STATUS_CH=$(grep -e ${CERT_NAME}$ -e${CERT_NAME}/ ${INDEX} | awk '{print $1}' | tr -d '\n')
if [[ $STATUS_CH = "V" ]]; then
    echo "Cert is VALID"
    echo "Will remove: ${CERT_SERIAL}"
else
    echo "Cert is REVOKED"
    echo "Will remove: ${CERT_SERIAL}"
fi

# Check if the user has two certificates in index.txt
if [[ $(cat /usr/share/easy-rsa/pki/index.txt | grep -c "/name=$NAME") -eq 2 ]]; then
    echo "Removing renewed certificate..."
    sed -i'.bak' "/${CERT_SERIAL}/d" $INDEX
    # removing *.ovpn file because it has old certificate
    rm -f $DEST_FILE_PATH
    
    echo 'Generate New .ovpn file...'
    CA="$(cat $EASY_RSA/pki/ca.crt )"
    CERT="$(cat $EASY_RSA/pki/issued/${CERT_NAME}.crt | grep -zEo -e '-----BEGIN CERTIFICATE-----(\n|.)*-----END CERTIFICATE-----' | tr -d '\0')"
    KEY="$(cat $EASY_RSA/pki/private/${CERT_NAME}.key)"
    TLS_AUTH="$(cat $EASY_RSA/pki/ta.key)"
    echo "$(cat /etc/openvpn/config/client.conf)
<ca>
$CA
</ca>
<cert>
$CERT
</cert>
<key>
$KEY
</key>
<tls-auth>
$TLS_AUTH
</tls-auth>
" > "$DEST_FILE_PATH"
    echo "all done"

else
    echo "Removing certificate..."
    rm -f /etc/openvpn/pki/certs_by_serial/$CERT_SERIAL.pem
    rm -f /etc/openvpn/pki/issued/$CERT_NAME.crt
    rm -f /etc/openvpn/pki/private/$CERT_NAME.key
    rm -f /etc/openvpn/pki/reqs/$CERT_NAME.req
    echo "removing *.ovpn file" 
    rm -f /etc/openvpn/clients/$CERT_NAME.ovpn

    # Fix index.txt by removing the user from the list following the serial number
    sed -i'.bak' "/${CERT_SERIAL}/d" $INDEX
fi

echo 'Remove done!'
echo 'If you want to disconnect the user please restart the OpenVPN service or container.'
