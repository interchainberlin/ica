#!/bin/bash
set -e

CONFIG_DIR=./network/hermes/config.toml

### Configure clients
echo "Configuring clients..."
hermes -c $CONFIG_DIR tx raw create-client test-1 test-2
hermes -c $CONFIG_DIR tx raw create-client test-2 test-1

### Connection Handshake
echo "Initiating connection handshake..."
# conn-init
hermes -c $CONFIG_DIR tx raw conn-init test-1 test-2 07-tendermint-0 07-tendermint-0
# conn-try
hermes -c $CONFIG_DIR tx raw conn-try test-2 test-1 07-tendermint-0 07-tendermint-0 -s connection-0
# conn-ack
hermes -c $CONFIG_DIR tx raw conn-ack test-1 test-2 07-tendermint-0 07-tendermint-0 -d connection-0 -s connection-0
# conn-confirm
hermes -c $CONFIG_DIR tx raw conn-confirm test-2 test-1 07-tendermint-0 07-tendermint-0 -d connection-0 -s connection-0


