#!/bin/bash
#VERSION 1.2 by d3vilh@github.com aka Mr. Philipp
# Exit immediately if a command exits with a non-zero status
set -e

#Variables
CERT_NAME=$1
CERT_SERIAL=$2
EASY_RSA=$(grep -E "^EasyRsaPath\s*=" ../openvpn-ui/conf/app.conf | cut -d= -f2 | tr -d '"' | tr -d '[:space:]')
OPENVPN_DIR=$(grep -E "^OpenVpnPath\s*=" ../openvpn-ui/conf/app.conf | cut -d= -f2 | tr -d '"' | tr -d '[:space:]')
echo 'EasyRSA path: $EASY_RSA OVPN path: $OPENVPN_DIR'
INDEX=$EASY_RSA/pki/index.txt
PERSHIY=`cat $INDEX | grep "/CN=$CERT_NAME/" | head -1 | awk '{ print $3}'`
OVPN_FILE_PATH="$OPENVPN_DIR/clients/$CERT_NAME.ovpn"

export EASYRSA_BATCH=1 # see https://superuser.com/questions/1331293/easy-rsa-v3-execute-build-ca-and-gen-req-silently

# Check if the user has two certificates in index.txt
if [[ $(cat $INDEX | grep -c "/CN=$CERT_NAME/") -eq 2 ]]; then
    # Check if first serial is the same as requested to revoke and if yes - revoke new cert and old cert
    if [[ $PERSHIY = $CERT_SERIAL ]]; then
        echo "Revoking renewed certificate..."

        # removing the end of the line starting from /name=$NAME for the line that matches the $serial pattern
        sed  -i'.bak' "/$CERT_SERIAL/s/\/name=$CERT_NAME.*//" $INDEX
        echo "index.txt patched"
        cd $EASY_RSA
        
        #moving new cert to old dir
        echo "Runing: easyrsa revoke-renewed $CERT_NAME"
        # Revoke renewed certificate
        ./easyrsa revoke-renewed "$CERT_NAME"
        echo -e "Old certificate revoked! \nRemoving old cert from the DB"

        # Removing old cert from the DB
        sed -i'.bak' "/${CERT_SERIAL}/d" $INDEX
        echo "Old cert with serial $CERT_SERIAL removed from the DB"

        # removing *.ovpn file because it has old certificate
        echo "removing *.ovpn file"
        rm -f $OVPN_FILE_PATH
    
        echo 'Generate New .ovpn file...'
        CA="$(cat $EASY_RSA/pki/ca.crt )"
        CERT="$(cat $EASY_RSA/pki/issued/${CERT_NAME}.crt | grep -zEo -e '-----BEGIN CERTIFICATE-----(\n|.)*-----END CERTIFICATE-----' | tr -d '\0')"
        KEY="$(cat $EASY_RSA/pki/private/${CERT_NAME}.key)"
        TLS_AUTH="$(cat $EASY_RSA/pki/ta.key)"
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
        echo -e "Old Certificate revoked!\nCreate new Create certificate revocation list (CRL)..."
        ./easyrsa gen-crl
        chmod +r $EASY_RSA/pki/crl.pem
    else
        cd $EASY_RSA
        # Fix index.txt by removing the user from the list following the serial number
        echo "Removing New Certificate..."
        mv $EASY_RSA/pki/renewed/issued/$CERT_NAME.crt  $EASY_RSA/pki/issued/$CERT_NAME.crt
        rm -f $EASY_RSA/pki/inline/$CERT_NAME.inline
        # Removing old cert from the DB
        sed -i'.bak' "/${CERT_SERIAL}/d" $INDEX
        # Create new Create certificate revocation list (CRL)
        echo -e "New Certificate revoked!\nCreate new certificate revocation list (CRL)..."
        ./easyrsa gen-crl
        chmod +r $EASY_RSA/pki/crl.pem
    fi
else
    echo "Revoking certificate..."
    # removing the end of the line starting from /name=$NAME for the line that matches the $serial pattern
    sed  -i'.bak' "/$CERT_SERIAL/s/\/name=$CERT_NAME.*//" $INDEX
    cd $EASY_RSA
    # Revoke certificate
    ./easyrsa revoke "$CERT_NAME"

    echo 'Create new Create certificate revocation list (CRL)...'
    ./easyrsa gen-crl
    chmod +r $EASY_RSA/pki/crl.pem
    # restoring the index.txt, new /name in index.txt (adding name and ip to the last line)
    sed -i'.bak' "$ s/$/\/name=${CERT_NAME}\/LocalIP=${CERT_IP}\/2FAName=${TFA_NAME}/" $EASY_RSA/pki/index.txt
fi

echo -e 'Done!\nIf you want to disconnect the user please restart the service using docker-compose restart openvpn.'
