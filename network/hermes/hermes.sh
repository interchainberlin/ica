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

### Create an ics-27 ibcaccount channel
echo "Creating ics-27 ibcaccount channel..."
# chan-open-init
hermes -c $CONFIG_DIR tx raw chan-open-init test-1 test-2 connection-0 ibcaccount ibcaccount -o ORDERED
# chan-open-try
hermes -c $CONFIG_DIR tx raw chan-open-try test-2 test-1 connection-0 ibcaccount ibcaccount -s channel-0
# chan-open-ack
hermes -c $CONFIG_DIR tx raw chan-open-ack test-1 test-2 connection-0 ibcaccount ibcaccount -d channel-0 -s channel-0
# chan-open-confirm
hermes -c $CONFIG_DIR tx raw chan-open-confirm test-2 test-1 connection-0 ibcaccount ibcaccount -d channel-0 -s channel-0

### Create an ics-20 transfer channel
echo "Creating ics-20 transfer channel..."
# chan-open-init
hermes -c $CONFIG_DIR tx raw chan-open-init test-1 test-2 connection-0 transfer transfer -o UNORDERED
# chan-open-try
hermes -c $CONFIG_DIR tx raw chan-open-try test-2 test-1 connection-0 transfer transfer -s channel-1
# chan-open-ack
hermes -c $CONFIG_DIR tx raw chan-open-ack test-1 test-2 connection-0 transfer transfer -d channel-1 -s channel-1
# chan-open-confirm
hermes -c $CONFIG_DIR tx raw chan-open-confirm test-2 test-1 connection-0 transfer transfer -d channel-1 -s channel-1

