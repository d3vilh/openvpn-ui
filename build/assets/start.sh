#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

# Directory where OpenVPN configuration files are stored
OVDIR=/etc/openvpn

# Change to the /opt directory
cd /opt/

# If the provisioned file does not exist in the OpenVPN directory, prepare the certificates and create the provisioned file
if [ ! -f $OVDIR/.provisioned ]; then
  #echo "Preparing certificates"
  mkdir -p $OVDIR

  # Uncomment line below to generate CA and server certificates (should be done on the side of OpenVPN container or server however)
  #./scripts/generate_ca_and_server_certs.sh

  # Create the provisioned file
  touch $OVDIR/.provisioned
  echo "First OpenVPN UI start."
fi

# Change to the OpenVPN GUI directory
cd /opt/openvpn-gui

# Create the database directory if it does not exist
mkdir -p db

# Start the OpenVPN GUI
./openvpn-ui
echo "Starting!"
