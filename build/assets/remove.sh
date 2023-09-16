#!/bin/bash
#VERSION 1.2 by @d3vilh@github.com aka Mr. Philipp
# Exit immediately if a command exits with a non-zero status
set -e

# Variables
ACTION=$1  #passed via OpenVPN-UI GUI
EASY_RSA=$(grep -E "^EasyRsaPath\s*=" ../openvpn-gui/conf/app.conf | cut -d= -f2 | tr -d '"' | tr -d '[:space:]')
OPENVPN_DIR=$(grep -E "^OpenVpnPath\s*=" ../openvpn-gui/conf/app.conf | cut -d= -f2 | tr -d '"' | tr -d '[:space:]')
echo 'EasyRSA path: $EASY_RSA OVPN path: $OPENVPN_DIR'
PKI_DIR=$EASY_RSA/pki
CERT_DIR=$OPENVPN_DIR/clients
STATIC_CLIENT_DIR=$OPENVPN_DIR/staticclients
OVPN_DB_DIR=$OPENVPN_DIR/db

if [ "$ACTION" = "remove_pki" ]; then
  # Remove PKI
  echo 'Removing PKI.'
  rm -rf $PKI_DIR/*
elif [ "$ACTION" = "remove_ovpn" ]; then
  # Remove *.ovpn files from /etc/openvpn/clients
  echo -e 'Removing *.ovpn files.\n'
  rm -rf $CERT_DIR/*.ovpn
elif [ "$ACTION" = "remove_static_clients" ]; then
  # Remove static clients
  echo -e 'Removing static clients.\n'
  rm -rf $STATIC_CLIENT_DIR/*
elif [ "$ACTION" = "remove_ovpn_db" ]; then
  # Remove Openvpn-UI DB
  echo -e 'Removing Openvpn-UI DB.\n'
  rm -rf $OVPN_DB_DIR/data.db
elif [ "$ACTION" = "remove_all" ]; then
  # Remove all
  echo 'Removing PKI.'
  rm -rf $PKI_DIR/*
  echo -e 'PKI removed.\nRemoving *.ovpn files.\n'
  rm -rf $CERT_DIR/*.ovpn
  echo -e 'All *.ovpn removed.\nRemoving static clients.\n'
  rm -rf $STATIC_CLIENT_DIR/*
  echo -e 'Static clients remover.\nRemoving Openvpn-UI DB.\n'
  rm -rf $OVPN_DB_DIR/data.db
else
  echo "Invalid input argument: $ACTION. Exiting."
  exit 1
fi