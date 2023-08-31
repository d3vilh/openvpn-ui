#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

# Variables
PKI_DIR=/usr/share/easy-rsa/pki
CERT_DIR=/etc/openvpn/clients
STATIC_CLIENT_DIR=/etc/openvpn/staticclients
OVPN_DB_DIR=/etc/openvpn/db

# Remove PKI
echo 'Removing PKI.'
rm -rf $PKI_DIR/*

# Remove certificates
echo -e 'PKI removed.\nRemoving certificates.\n'
rm -rf $CERT_DIR/*.ovpn

# Remove static clients
echo -e 'Certificates removed.\nRemoving static clients.\n'
rm -rf $STATIC_CLIENT_DIR/*

# Remove Openvpn-UI DB
echo -e 'Static clients removed.\nRemoving Openvpn-UI DB.\n'
rm -rf $OVPN_DB_DIR/data.db

echo 'Done. PKI, Certificates, Satic Clients and OpenVPN-UI DB removed.' 