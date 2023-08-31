#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

# Variables
PKI_DIR=/usr/share/easy-rsa/pki
CERT_DIR=/etc/openvpn/clients
STATIC_CLIENT_DIR=/etc/openvpn/staticclients
OVPN_DB_DIR=/etc/openvpn/db

if [ "$1" = "remove_pki" ]; then
  # Remove PKI
  echo 'Removing PKI.'
  rm -rf $PKI_DIR/*
elif [ "$1" = "remove_cert" ]; then
  # Remove certificates
  echo -e 'PKI removed.\nRemoving certificates.\n'
  rm -rf $CERT_DIR/*.ovpn
elif [ "$1" = "remove_static_clients" ]; then
  # Remove static clients
  echo -e 'Certificates removed.\nRemoving static clients.\n'
  rm -rf $STATIC_CLIENT_DIR/*
elif [ "$1" = "remove_ovpn_db" ]; then
  # Remove Openvpn-UI DB
  echo -e 'Static clients removed.\nRemoving Openvpn-UI DB.\n'
  rm -rf $OVPN_DB_DIR/data.db
elif [ "$1" = "remove_all" ]; then
  # Remove all
  echo 'Removing PKI.'
  rm -rf $PKI_DIR/*
  echo -e 'PKI removed.\nRemoving certificates.\n'
  rm -rf $CERT_DIR/*.ovpn
  echo -e 'Certificates removed.\nRemoving static clients.\n'
  rm -rf $STATIC_CLIENT_DIR/*
  echo -e 'Static clients remover.\nRemoving Openvpn-UI DB.\n'
  rm -rf $OVPN_DB_DIR/data.db
else
  echo "Invalid input argument: $1"
  exit 1
fi