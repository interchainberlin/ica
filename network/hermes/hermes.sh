#!/bin/bash
set -e

CONFIG_DIR=./network/hermes/config.toml

### Configure clients
echo "Configuring clients..."
hermes -c $CONFIG_DIR tx raw create-client test-1 test-2
hermes -c $CONFIG_DIR tx raw create-client test-2 test-1

### Connection Handshake
echo "Initiating connection handshake..."
# create clients, connection and channel
hermes -c $CONFIG_DIR create connection test-1 test-2
sleep 2


