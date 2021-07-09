#!/bin/bash

# Start the hermes relayer in multi-paths mode
CONFIG_DIR=./network/hermes/config.toml
PORT_ID=ics27-1-0-cosmos1mjk79fjjgpplak5wq838w0yd982gzkyfrk07am
hermes -c ./network/hermes/config.toml tx raw chan-open-init test-1 test-2 connection-0 $PORT_ID ibcaccount -o ORDERED
hermes -c $CONFIG_DIR tx raw chan-open-try test-2 test-1 connection-0 ibcaccount $PORT_ID -s channel-0
hermes -c $CONFIG_DIR tx raw chan-open-ack test-1 test-2 connection-0 $PORT_ID ibcaccount -d channel-0 -s channel-0
hermes -c $CONFIG_DIR tx raw chan-open-confirm test-2 test-1 connection-0 ibcaccount $PORT_ID -d channel-0 -s channel-0


