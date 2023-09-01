#!/bin/bash -e
EASY_RSA=/usr/share/easy-rsa
OPENVPN_DIR=/etc/openvpn
TEMP_PKI_DIR=/tmp/pki
mkdir -p $TEMP_PKI_DIR

cd $EASY_RSA
 
if [ "$1" = "copy_vars" ]; then
  # Copy easy-rsa variables
  cp $OPENVPN_DIR/config/easy-rsa.vars ./pki/vars
  echo 'New vars applied.'
fi

if [[ ! -f $OPENVPN_DIR/pki/openssl-easyrsa.cnf || ! -f $OPENVPN_DIR/pki/ca.crt || ! -f $OPENVPN_DIR/pki/issued/server.crt || ! -f $OPENVPN_DIR/pki/dh.pem || ! -f $OPENVPN_DIR/pki/ta.key || ! -f $OPENVPN_DIR/pki/crl.pem || "$1" = "init_all" || "$1" = "gen_crl" ]] && ! [[ "$1" = "copy_vars" ]]; then
    export EASYRSA_BATCH=1 # see https://superuser.com/questions/1331293/easy-rsa-v3-execute-build-ca-and-gen-req-silently

    # Copy easy-rsa variables
    cp $OPENVPN_DIR/config/easy-rsa.vars ./pki/vars

    echo "Following EASYRSA variables will be used:"
    cat $EASY_RSA/pki/vars | awk '{$1=""; print $0}';

  if [[ "$1" = "init-pki" && ! -f $OPENVPN_DIR/pki/openssl-easyrsa.cnf ]]; then
      # Building the CA with WA to avoid issues with .pki container volume which not possible to remove due to its origin
      echo 'Setting up public key infrastructure...'
      $EASY_RSA/easyrsa --pki-dir=$TEMP_PKI_DIR init-pki

      echo 'Moving PKI directory...'
      mv $TEMP_PKI_DIR/* ./pki/
      cp $OPENVPN_DIR/config/easy-rsa.vars ./pki/vars

    elif [[ "$1" = "build_ca" && ! -f $OPENVPN_DIR/pki/ca.crt ]]; then
      echo 'Generating Certificate authority...'
      $EASY_RSA/easyrsa build-ca nopass

    elif [[ "$1" = "gen_req" && ! -f $OPENVPN_DIR/pki/issued/server.crt ]]; then
      # Creating the Server Certificate, Key, and Encryption Files
      echo 'Creating the Server Certificate...'
      $EASY_RSA/easyrsa gen-req server nopass

      echo 'Sign request...'
      $EASY_RSA/easyrsa sign-req server server

    elif [[ "$1" = "gen_dh" && ! -f $OPENVPN_DIR/pki/dh.pem ]]; then
      echo 'Generate Diffie-Hellman key...'
      $EASY_RSA/easyrsa gen-dh

    elif [[ "$1" = "gen_ta" && ! -f $OPENVPN_DIR/pki/ta.key ]]; then
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

    elif [ "$1" = "gen_crl" ]; then
      echo 'Create certificate revocation list (CRL)...'
      $EASY_RSA/easyrsa gen-crl
      chmod +r ./pki/crl.pem

    elif [ "$1" = "init_all" ]; then
      # Init all Begin
      cp $OPENVPN_DIR/config/easy-rsa.vars ./pki/vars

      echo "Following EASYRSA variables will be used:"
      cat $EASY_RSA/pki/vars | awk '{$1=""; print $0}';

      echo 'Setting up public key infrastructure...'
      $EASY_RSA/easyrsa --pki-dir=$TEMP_PKI_DIR init-pki

      echo 'Moving PKI directory...'
      mv $TEMP_PKI_DIR/* ./pki/

      cp $OPENVPN_DIR/config/easy-rsa.vars ./pki/vars

      echo 'Generating Certificate authority...'
      $EASY_RSA/easyrsa build-ca nopass

      echo 'Creating the Server Certificate...'
      $EASY_RSA/easyrsa gen-req server nopass

      echo 'Sign request...'
      $EASY_RSA/easyrsa sign-req server server

      echo 'Generate Diffie-Hellman key...'
      $EASY_RSA/easyrsa gen-dh

      echo 'Generate HMAC signature...'
      if ! command -v openvpn &> /dev/null
      then 
        echo 'Running in Docker container...'
        CONTAINER_ID=$(curl --unix-socket /var/run/docker.sock "http://v1.40/containers/json?filters=%7B%22name%22%3A%5B%22%5Eopenvpn$%22%5D%7D" | jq -r '.[0].Id')
        EXEC_ID=$(curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" -d '{"AttachStdin": true, "AttachStdout": true, "AttachStderr": true, "Cmd": ["openvpn", "--genkey", "--secret", "'"$OPENVPN_DIR"'/pki/ta.key"], "DetachKeys": "ctrl-p,ctrl-q", "Tty": true}' -X POST "http://v1.40/containers/$CONTAINER_ID/exec" | jq -r '.Id')
        curl --unix-socket /var/run/docker.sock -H "Content-Type: application/json" -d '{"Detach": false, "Tty": true}' -X POST "http://v1.40/exec/$EXEC_ID/start"
      else
         echo 'Running in host...'
         openvpn --genkey --secret $OPENVPN_DIR/pki/ta.key
      fi

      echo 'Create certificate revocation list (CRL)...'
      $EASY_RSA/easyrsa gen-crl
      chmod +r ./pki/crl.pem
      # Init all End
    else
      echo "Invalid input argument: $1"
      exit 1
  fi

   echo 'All done.'
else
    echo 'PKI already set up.'
fi
