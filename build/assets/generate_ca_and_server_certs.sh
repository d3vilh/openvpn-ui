#!/bin/bash -e

EASY_RSA=/usr/share/easy-rsa

if [[ ! -f /etc/openvpn/pki/ca.crt ]]; then
    export EASYRSA_BATCH=1 # see https://superuser.com/questions/1331293/easy-rsa-v3-execute-build-ca-and-gen-req-silently
    cd $EASY_RSA
 
    # Copy easy-rsa variables
    cp /etc/openvpn/config/easy-rsa.vars ./vars

    # Listing env parameters:
    echo "Following EASYRSA variables will be used:"
    cat $EASY_RSA/vars | awk '{$1=""; print $0}';

    # Building the CA
    echo 'Setting up public key infrastructure...'
    $EASY_RSA/easyrsa init-pki

    echo 'Generating ertificate authority...'
    $EASY_RSA/easyrsa build-ca nopass

    # Creating the Server Certificate, Key, and Encryption Files
    echo 'Creating the Server Certificate...'
    $EASY_RSA/easyrsa gen-req server nopass

    echo 'Sign request...'
    $EASY_RSA/easyrsa sign-req server server

    echo 'Generate Diffie-Hellman key...'
    $EASY_RSA/easyrsa gen-dh

    echo 'Generate HMAC signature...'
    openvpn --genkey --secret pki/ta.key

    echo 'Create certificate revocation list (CRL)...'
    $EASY_RSA/easyrsa gen-crl
    chmod +r ./pki/crl.pem

    # Copy to mounted volume
    cp -r ./pki/. /etc/openvpn/pki
else
    # Copy from mounted volume
    cp -r /etc/openvpn/pki /opt/app/easy-rsa
    echo 'PKI already set up.'
fi
