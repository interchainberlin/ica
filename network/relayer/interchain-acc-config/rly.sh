#!/bin/bash

# Configure predefined mnemonic pharses
BINARY=rly
CHAINDIR=./data
RELAYERDIR=./relayer
MNEMONIC1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
MNEMONIC2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"

# Ensure rly is installed
if ! [ -x "$(command -v $BINARY)" ]; then
    echo "$BINARY is required to run this script..."
    echo "You can download at https://github.com/cosmos/relayer"
    exit 1
fi

echo "Initializing $BINARY..."
$BINARY config init --home $CHAINDIR/$RELAYERDIR

echo "Adding configurations for both chains..."
$BINARY config add-chains $PWD/network/relayer/interchain-acc-config/chains --home $CHAINDIR/$RELAYERDIR
$BINARY config add-paths $PWD/network/relayer/interchain-acc-config/paths --home $CHAINDIR/$RELAYERDIR

echo "Restoring accounts..."
$BINARY keys restore test-1 test-1 "$MNEMONIC1" --home $CHAINDIR/$RELAYERDIR
$BINARY keys restore test-2 test-2 "$MNEMONIC2" --home $CHAINDIR/$RELAYERDIR

echo "Initializing light clients for both chains..."
$BINARY light init test-1 -f --home $CHAINDIR/$RELAYERDIR
$BINARY light init test-2 -f --home $CHAINDIR/$RELAYERDIR

echo "Linking both chains..."
$BINARY tx path test1-account-test2 --home $CHAINDIR/$RELAYERDIR

echo "Starting to listen relayer..."
$BINARY start test1-account-test2 --home $CHAINDIR/$RELAYERDIR
