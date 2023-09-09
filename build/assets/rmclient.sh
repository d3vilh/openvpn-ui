#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

NAME=$1
SERIAL=$2
INDEX=/usr/share/easy-rsa/pki/index.txt
EASY_RSA="/usr/share/easy-rsa"
PERSHIY=`cat $INDEX | grep "/name=$NAME" | head -1 | awk '{ print $3}'`
#ACTION=$3

# .ovpn file path
#DEST_FILE_PATH="/etc/openvpn/clients/$NAME.ovpn"
# Check if .ovpn file exists
#if [[ ! -f $DEST_FILE_PATH ]]; then
#    echo "User not found."
#    exit 1
#fi

export EASYRSA_BATCH=1 # see https://superuser.com/questions/1331293/easy-rsa-v3-execute-build-ca-and-gen-req-silently

# Check if the user has two certificates in index.txt
if [[ $(cat $INDEX | grep -c "/name=$NAME") -eq 2 ]]; then
    # Check if first serial is the same as requested to revoke and if yes - revoke new cert and old cert
    if [[ $PERSHIY = $SERIAL ]]; then
        echo "Revoking renewed certificate and then old one..."
        # Fix index.txt by removing everything after pattern "/name=$1" in the line
        sed -i'.bak' "s/\/name=${NAME}\/.*//" $INDEX
        cd $EASY_RSA
        # Revoke renewed certificate
        ./easyrsa revoke-renewed "$NAME"
        echo -e "Renewed certificate revoked!/nRevoking old certificate..."
        cd $EASY_RSA
        # Revoke certificate
        ./easyrsa revoke "$NAME"
        # Create new Create certificate revocation list (CRL)
        echo -e "Old certificate revoked!/nCreate new Create certificate revocation list (CRL)..."
        ./easyrsa gen-crl
        chmod +r ./pki/crl.pem
    else
        echo "Revoking renewed certificate..."
        # removing the end of the line starting from /name=$NAME for the line that matches the $serial pattern
        sed  -i'.bak' "/$SERIAL/s/\/name=$NAME.*//" $INDEX
        cd $EASY_RSA
        # Revoke renewed certificate
        ./easyrsa revoke-renewed "$NAME"
    fi
else
    echo "Revoking certificate..."
    # removing the end of the line starting from /name=$NAME for the line that matches the $serial pattern
    sed  -i'.bak' "/$SERIAL/s/\/name=$NAME.*//" $INDEX
    cd $EASY_RSA
    # Revoke certificate
    ./easyrsa revoke "$NAME"

    echo 'Create new Create certificate revocation list (CRL)...'
    ./easyrsa gen-crl
    chmod +r ./pki/crl.pem
fi

echo 'Done!'
echo 'If you want to disconnect the user please restart the service using docker-compose restart openvpn.'


# Old Revoke:
# Fix index.txt by removing everything after pattern "/name=$1" in the line
# Fix for https://github.com/d3vilh/openvpn-ui/issues/5 by shuricksumy@github
#sed -i'.bak' "s/\/name=${name}\/.*//" /usr/share/easy-rsa/pki/index.txt
