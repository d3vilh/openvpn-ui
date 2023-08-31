#!/bin/bash -e
EASY_RSA=/usr/share/easy-rsa
PKI_DIR=/tmp/pki
mkdir -p $PKI_DIR

if [[ ! -f /etc/openvpn/pki/ca.crt ]]; then
    export EASYRSA_BATCH=1 # see https://superuser.com/questions/1331293/easy-rsa-v3-execute-build-ca-and-gen-req-silently
    cd $EASY_RSA
 
    # Copy easy-rsa variables
    cp /etc/openvpn/config/easy-rsa.vars ./pki/vars

    # Listing env parameters:
    echo "Following EASYRSA variables will be used:"
    cat $EASY_RSA/vars | awk '{$1=""; print $0}';

    # Building the CA
    echo 'Setting up public key infrastructure...'
    $EASY_RSA/easyrsa --pki-dir=$PKI_DIR init-pki

    echo 'Moving PKI directory...'
    mv $PKI_DIR/* ./pki/

    echo 'Generating ertificate authority...'
    $EASY_RSA/easyrsa build-ca nopass

    # Creating the Server Certificate, Key, and Encryption Files
    echo 'Creating the Server Certificate...'
    $EASY_RSA/easyrsa gen-req server nopass

    echo 'Sign request...'
    $EASY_RSA/easyrsa sign-req server server

    echo 'Generate Diffie-Hellman key...'
    $EASY_RSA/easyrsa gen-dh

    # Generate HMAC signature in "openvpn" container with Docker API to strengthen server certificate.
    echo 'Generate HMAC signature...'
    nohup sh -c 'curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" -d '\''{  "Cmd": ["openvpn", "--genkey", "--secret", "/etc/openvpn/pki/ta.key"]}'\'' -X POST http://$DOCKER_HOST_IP/containers/$CONTAINER_ID/exec | jq -r '\''{.Id}'\'' | xargs curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" -d '\''{"Detach": false, "Tty": false}'\'' -X POST http://$DOCKER_HOST_IP/exec/{}/start' > /dev/null 2>&1 &

    # curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" -d '{  "Cmd": ["openvpn", "--genkey", "--secret", "/etc/openvpn/pki/ta.key"] }' -X POST http://$DOCKER_HOST_IP/containers/$CONTAINER_ID/exec | jq -r '.Id' | xargs curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" -d '{"Detach": false, "Tty": false}' -X POST http://$DOCKER_HOST_IP/exec/{}/start

    echo 'Create certificate revocation list (CRL)...'
    $EASY_RSA/easyrsa gen-crl
    chmod +r ./pki/crl.pem

    echo 'All done.'
    # Copy to mounted volume
    # cp -r ./pki/. /etc/openvpn/pki
else
    # Copy from mounted volume
    # cp -r /etc/openvpn/pki /opt/app/easy-rsa
    echo 'PKI already set up.'
fi
