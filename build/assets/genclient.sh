#!/bin/bash
#VERSION 1.2 by @d3vilh@github.com aka Mr. Philipp
# Exit immediately if a command exits with a non-zero status
set -e

# .ovpn file path
CERT_NAME=$1
CERT_IP=$2
CERT_PASS=$3
#EASY_RSA=/usr/share/easy-rsa
#OPENVPN_DIR=/etc/openvpn
EASY_RSA=$(grep -E "^EasyRsaPath\s*=" ../openvpn-gui/conf/app.conf | cut -d= -f2 | tr -d '[:space:]')
OPENVPN_DIR=$(grep -E "^OpenVpnPath\s*=" ../openvpn-gui/conf/app.conf | cut -d= -f2 | tr -d '[:space:]')
echo 'EasyRSA path: $EASY_RSA OVPN path: $OPENVPN_DIR'
OVPN_FILE_PATH="$OPENVPN_DIR/clients/$CERT_NAME.ovpn"

# Validate username and check for duplicates
if  [[ -z $CERT_NAME ]]; then
    echo 'Name cannot be empty. Exiting...'
    exit 1
elif [[ -f $OVPN_FILE_PATH ]]; then
    echo "User with name $CERT_NAME already exists under openvpn/clients. Exiting..."
    exit 1
fi

export EASYRSA_BATCH=1 # see https://superuser.com/questions/1331293/easy-rsa-v3-execute-build-ca-and-gen-req-silently

echo 'Patching easy-rsa.3.1.1 openssl-easyrsa.cnf...' 
sed -i '/serialNumber_default/d' "$EASY_RSA/pki/openssl-easyrsa.cnf"

echo 'Generate client certificate...'

# Copy easy-rsa variables
cd $EASY_RSA

# Generate certificates
if  [[ -z $CERT_PASS ]]; then
    echo 'Without password...'
    ./easyrsa --batch --req-cn="$CERT_NAME" gen-req "$CERT_NAME" nopass 
else
    echo 'With password...'
    # See https://stackoverflow.com/questions/4294689/how-to-generate-an-openssl-key-using-a-passphrase-from-the-command-line
    # ... and https://stackoverflow.com/questions/22415601/using-easy-rsa-how-to-automate-client-server-creation-process
    # ... and https://github.com/OpenVPN/easy-rsa/blob/master/doc/EasyRSA-Advanced.md
    (echo -e '\n') | ./easyrsa --batch --req-cn="$CERT_NAME" --passin=pass:${CERT_PASS} --passout=pass:${CERT_PASS} gen-req "$CERT_NAME"
fi

# Sign request
./easyrsa sign-req client "$CERT_NAME"
# Fix for /name in index.txt
echo "Fixind Database..."
sed -i'.bak' "$ s/$/\/name=${CERT_NAME}\/LocalIP=${CERT_IP}/" $EASY_RSA/pki/index.txt
# Certificate properties
CA="$(cat $EASY_RSA/pki/ca.crt )"
CERT="$(cat $EASY_RSA/pki/issued/${CERT_NAME}.crt | grep -zEo -e '-----BEGIN CERTIFICATE-----(\n|.)*-----END CERTIFICATE-----' | tr -d '\0')"
KEY="$(cat $EASY_RSA/pki/private/${CERT_NAME}.key)"
TLS_AUTH="$(cat $EASY_RSA/pki/ta.key)"

echo 'Fixing permissions for pki/issued...'
chmod +r $EASY_RSA/pki/issued

echo 'Generating .ovpn file...'
echo "$(cat $OPENVPN_DIR/config/client.conf)
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
" > "$OVPN_FILE_PATH"

echo "OpenVPN Client configuration successfully generated!\nCheckout openvpn/clients/$CERT_NAME.ovpn"
