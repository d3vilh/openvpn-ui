#!/bin/sh
# v.0.1 by @d3vilh@github.com aka Mr. Philipp
# d3vilh/openvpn-server drafted 2FA support
#
# MFA verification by OpenVPN server using oath-tool
# This should be part of OpenVPN server container (github.com/d3vilh/openvpn-server). 
# This particular script is just example and does not use by OpenVPN-UI at all.

# VARIABLES
PASSFILE=$1    # Password file passed by openvpn-server with "auth-user-pass-verify /opt/app/bin/oath.sh via-file" in server.conf
OPENVPN_DIR=/etc/openvpn
OATH_SECRETS=$OPENVPN_DIR/clients/oath.secrets
LOG_FILE=/var/log/openvpn/oath.log

echo -e "$(date) Openvpn dir: $OPENVPN_DIR\nOath secrets: $OATH_SECRETS\nLog file: $LOG_FILE\nPassfile: $PASSFILE\n" | tee -a $LOG_FILE

# Geting user and password passed by external user to OpenVPN server tmp file
user=$(head -1 $PASSFILE)
pass=$(tail -1 $PASSFILE) 

echo "$(date) - Authentication attempt for user $user" | tee -a $LOG_FILE # echo "$(date) - Password: $pass" | tee -a $LOG_FILE

# Parsing oath.secrets to getting secret entry, ignore case
secret=$(grep -i -m 1 "$user:" $OATH_SECRETS | cut -d: -f2) # echo "$(date) - Secret: $secret" | tee -a $LOG_FILE

# Getting 2FA code with oathtool based on our secret, exiting with 0 if match:
code=$(oathtool --totp $secret) # echo "$(date) - Code: $code" | tee -a $LOG_FILE

if [ "$code" = "$pass" ];
then
    echo "OK"
        exit 0
else 
echo "FAIL"
fi

# If we make it here, auth hasn't succeeded, don't grant access
echo "$(date) - Authentication failed for user $user" | tee -a $LOG_FILE
exit 1