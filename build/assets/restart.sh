#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

# Get the container ID for ^openvpn$
CONTAINER_ID=$(curl --unix-socket /var/run/docker.sock "http://v1.40/containers/json?filters=%7B%22name%22%3A%5B%22%5Eopenvpn$%22%5D%7D" | grep '"Id":' | cut -d '"' -f 4)

# Restart the container
curl --unix-socket /var/run/docker.sock -X POST "http://v1.40/containers/$CONTAINER_ID/restart"

# Print the restarted container ID
echo "Restarted container $CONTAINER_ID"
