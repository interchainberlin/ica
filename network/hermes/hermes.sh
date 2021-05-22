#!/bin/bash
set -e

### Configure clients
echo "Configuring clients..."
hermes -c config.toml tx raw create-client test-1 test-2
hermes -c config.toml tx raw create-client test-2 test-1

### Connection Handshake
echo "Initiating connection handshake..."
# conn-init
hermes -c config.toml tx raw conn-init test-1 test-2 07-tendermint-0 07-tendermint-0
# conn-try
hermes -c config.toml tx raw conn-try test-2 test-1 07-tendermint-0 07-tendermint-0 -s connection-0
# conn-ack
hermes -c config.toml tx raw conn-ack test-1 test-2 07-tendermint-0 07-tendermint-0 -d connection-0 -s connection-0
# conn-confirm
hermes -c config.toml tx raw conn-confirm test-2 test-1 07-tendermint-0 07-tendermint-0 -d connection-0 -s connection-0

### Create an ics-20 transfer channel
echo "Creating ics-20 transfer channel..."
# chan-open-init
hermes -c config.toml tx raw chan-open-init test-1 test-2 connection-0 transfer transfer -o UNORDERED
# chan-open-try
hermes -c config.toml tx raw chan-open-try test-2 test-1 connection-0 transfer transfer -s channel-0
# chan-open-ack
hermes -c config.toml tx raw chan-open-ack test-1 test-2 connection-0 transfer transfer -d channel-0 -s channel-0
# chan-open-confirm
hermes -c config.toml tx raw chan-open-confirm test-2 test-1 connection-0 transfer transfer -d channel-0 -s channel-0

### Create an ics-27 ibcaccount channel
echo "Creating ics-27 ibcaccount channel..."
# chan-open-init
hermes -c config.toml tx raw chan-open-init test-1 test-2 connection-0 ibcaccount ibcaccount -o ORDERED
# chan-open-try
hermes -c config.toml tx raw chan-open-try test-2 test-1 connection-0 ibcaccount ibcaccount -s channel-1
# chan-open-ack
hermes -c config.toml tx raw chan-open-ack test-1 test-2 connection-0 ibcaccount ibcaccount -d channel-1 -s channel-1
# chan-open-confirm
hermes -c config.toml tx raw chan-open-confirm test-2 test-1 connection-0 ibcaccount ibcaccount -d channel-1 -s channel-1

# Start the hermes relayer in multi-paths mode
echo "Starting hermes relayer..."
hermes -c config.toml start-multi