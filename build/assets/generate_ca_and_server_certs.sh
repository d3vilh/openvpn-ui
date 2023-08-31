#!/bin/bash -e
EASY_RSA=/usr/share/easy-rsa
OPENVPN_DIR=/etc/openvpn
TEMP_PKI_DIR=/tmp/pki
mkdir -p $TEMP_PKI_DIR

if [[ ! -f $OPENVPN_DIR/pki/ca.crt ]]; then
    export EASYRSA_BATCH=1 # see https://superuser.com/questions/1331293/easy-rsa-v3-execute-build-ca-and-gen-req-silently
    cd $EASY_RSA
 
    # Copy easy-rsa variables
    cp $OPENVPN_DIR/config/easy-rsa.vars ./pki/vars

    # Listing env parameters:
    echo "Following EASYRSA variables will be used:"
    cat $EASY_RSA/pki/vars | awk '{$1=""; print $0}';

    # Building the CA with WA to avoid issues with .pki container volume which not possible to remove due to its origin
    echo 'Setting up public key infrastructure...'
    $EASY_RSA/easyrsa --pki-dir=$TEMP_PKI_DIR init-pki

    echo 'Moving PKI directory...'
    mv $TEMP_PKI_DIR/* ./pki/

    echo 'Generating Certificate authority...'
    $EASY_RSA/easyrsa build-ca nopass

    # Creating the Server Certificate, Key, and Encryption Files
    echo 'Creating the Server Certificate...'
    $EASY_RSA/easyrsa gen-req server nopass

    echo 'Sign request...'
    $EASY_RSA/easyrsa sign-req server server

    echo 'Generate Diffie-Hellman key...'
    $EASY_RSA/easyrsa gen-dh

    # Generate HMAC signature in "openvpn" container with Docker API or in host
    echo 'Generate HMAC signature...'
    # Check if the "openvpn" command exists
    if ! command -v openvpn &> /dev/null
    then 
      echo 'Running in Docker container...'
      # Get the container ID of the "openvpn" container
      CONTAINER_ID=$(curl --unix-socket /var/run/docker.sock "http://v1.40/containers/json?filters=%7B%22name%22%3A%5B%22%5Eopenvpn$%22%5D%7D" | jq -r '.[0].Id')

      # Create the exec instance
      EXEC_ID=$(curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" -d '{"AttachStdin": true, "AttachStdout": true, "AttachStderr": true, "Cmd": ["openvpn", "--genkey", "--secret", "'"$OPENVPN_DIR"'/pki/ta.key"], "DetachKeys": "ctrl-p,ctrl-q", "Tty": true}' -X POST "http://v1.40/containers/$CONTAINER_ID/exec" | jq -r '.Id')

      # Start the exec instance
      curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" -d '{"Detach": false, "Tty": true}' -X POST "http://v1.40/exec/$EXEC_ID/start"
    else
       # Run the "openvpn --genkey --secret pki/ta.key" command on localhost
       echo 'Running in host...'
       openvpn --genkey --secret $OPENVPN_DIR/pki/ta.key
    fi
    echo 'Create certificate revocation list (CRL)...'
    $EASY_RSA/easyrsa gen-crl
    chmod +r ./pki/crl.pem

    echo 'All done.'
else
    echo 'PKI already set up.'
fi
