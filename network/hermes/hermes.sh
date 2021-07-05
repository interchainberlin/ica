#!/bin/bash
set -e

CONFIG_DIR=./network/hermes/config.toml

### Configure the clients and connection
echo "Initiating connection handshake..."
hermes -c $CONFIG_DIR create connection test-1 test-2

sleep 2


