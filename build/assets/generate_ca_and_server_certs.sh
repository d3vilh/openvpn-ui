#!/bin/bash -e
#VERSION 1.2 by d3vilh@github.com aka Mr. Philipp

#Variables
ACTION=$1  #passed via OpenVPN-UI GUI
EASY_RSA=$(grep -E "^EasyRsaPath\s*=" ../openvpn-ui/conf/app.conf | cut -d= -f2 | tr -d '"' | tr -d '[:space:]')
OPENVPN_DIR=$(grep -E "^OpenVpnPath\s*=" ../openvpn-ui/conf/app.conf | cut -d= -f2 | tr -d '"' | tr -d '[:space:]')
echo "EasyRSA path: $EASY_RSA OVPN path: $OPENVPN_DIR"
TEMP_PKI_DIR=/tmp/pki
mkdir -p $TEMP_PKI_DIR

cd $EASY_RSA
 
if [ "$ACTION" = "copy_vars" ]; then
  # Copy easy-rsa variables
  cp $OPENVPN_DIR/config/easy-rsa.vars $EASY_RSA/pki/vars
  echo 'New vars applied.'
fi

if [[ ! -f $OPENVPN_DIR/pki/openssl-easyrsa.cnf || ! -f $OPENVPN_DIR/pki/ca.crt || ! -f $OPENVPN_DIR/pki/issued/server.crt || ! -f $OPENVPN_DIR/pki/dh.pem || ! -f $OPENVPN_DIR/pki/ta.key || ! -f $OPENVPN_DIR/pki/crl.pem || "$ACTION" = "init_all" || "$ACTION" = "gen_crl" ]] && ! [[ "$ACTION" = "copy_vars" ]]; then
    export EASYRSA_BATCH=1 # see https://superuser.com/questions/1331293/easy-rsa-v3-execute-build-ca-and-gen-req-silently

    # Copy easy-rsa variables
    cp $OPENVPN_DIR/config/easy-rsa.vars $EASY_RSA/pki/vars

    echo "Following EASYRSA variables will be used:"
    cat $EASY_RSA/pki/vars | awk '{$1=""; print $0}';

  if [[ "$ACTION" = "init-pki" && ! -f $OPENVPN_DIR/pki/openssl-easyrsa.cnf ]]; then
      # Building the CA with WA to avoid issues with .pki container volume which not possible to remove due to its origin
      echo 'Setting up public key infrastructure...'
      $EASY_RSA/easyrsa --pki-dir=$TEMP_PKI_DIR init-pki

      echo 'Moving PKI directory...'
      mv $TEMP_PKI_DIR/* $EASY_RSA/pki/
      cp $OPENVPN_DIR/config/easy-rsa.vars $EASY_RSA/pki/vars

    elif [[ "$ACTION" = "build_ca" && ! -f $OPENVPN_DIR/pki/ca.crt ]]; then
      echo 'Generating Certificate authority...'
      $EASY_RSA/easyrsa build-ca nopass

    elif [[ "$ACTION" = "gen_req" && ! -f $OPENVPN_DIR/pki/issued/server.crt ]]; then
      # Creating the Server Certificate, Key, and Encryption Files
      echo 'Creating the Server Certificate...'
      $EASY_RSA/easyrsa gen-req server nopass

      echo 'Sign request...'
      $EASY_RSA/easyrsa sign-req server server

    elif [[ "$ACTION" = "gen_dh" && ! -f $OPENVPN_DIR/pki/dh.pem ]]; then
      echo 'Generate Diffie-Hellman key...'
      $EASY_RSA/easyrsa gen-dh

    elif [[ "$ACTION" = "gen_ta" && ! -f $OPENVPN_DIR/pki/ta.key ]]; then
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

    elif [ "$ACTION" = "gen_crl" ]; then
      echo 'Create certificate revocation list (CRL)...'
      $EASY_RSA/easyrsa gen-crl
      chmod +r $EASY_RSA/pki/crl.pem

    elif [ "$ACTION" = "init_all" ]; then
      # Init all Begin
      cp $OPENVPN_DIR/config/easy-rsa.vars $EASY_RSA/pki/vars

      echo "Following EASYRSA variables will be used:"
      cat $EASY_RSA/pki/vars | awk '{$1=""; print $0}';

      echo 'Setting up public key infrastructure...'
      $EASY_RSA/easyrsa --pki-dir=$TEMP_PKI_DIR init-pki

      echo 'Moving PKI directory...'
      mv $TEMP_PKI_DIR/* $EASY_RSA/pki/

      cp $OPENVPN_DIR/config/easy-rsa.vars $EASY_RSA/pki/vars

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
      chmod +r $EASY_RSA/pki/crl.pem
      # Init all End
    else
      echo "Invalid input argument: $ACTION Exiting."
      exit 1
  fi

   echo 'All done.'
else
    echo 'PKI already set up.'
fi
