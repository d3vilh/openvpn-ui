#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e
echo 'Removing PKI.'
rm -rf /usr/share/easy-rsa/pki/*
echo 'PKI removed.'
echo 'Removing certificates.'
rm -rf /etc/openvpn/clients/*.ovpn
echo 'Removing static clients.'
rm -rf /etc/openvpn/staticclients/*
echo 'Removing Openvpn-UI DB.'
rm -rf /etc/openvpn/db/data.db
echo 'Done' 