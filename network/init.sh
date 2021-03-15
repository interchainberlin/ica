#!/bin/bash

# Configure predetermined settings
# val1: cosmos1mjk79fjjgpplak5wq838w0yd982gzkyfrk07am
# val2: cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs
BINARY=icad
CHAINDIR=./data
CHAINID1=test-1
CHAINID2=test-2
MNEMONIC1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
MNEMONIC2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"
P2PPORT1=16656
P2PPORT2=26656
RPCPORT1=16657
RPCPORT2=26657
GRPCPORT1=8090
GRPCPORT2=9090
RESTPORT1=1316
RESTPORT2=1317

# Stop if it is already running 
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    killall $BINARY
fi

# Remove previous data
echo "Removing previous data..."
rm -rf $CHAINDIR/$CHAINID1 &> /dev/null
rm -rf $CHAINDIR/$CHAINID2 &> /dev/null

# Add directories for both chains, exit if an error occurs
if ! mkdir -p $CHAINDIR/$CHAINID1 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

if ! mkdir -p $CHAINDIR/$CHAINID2 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

echo "Initializing $CHAINID1..."
echo "Initializing $CHAINID2..."
$BINARY init test --home $CHAINDIR/$CHAINID1 --chain-id=$CHAINID1
$BINARY init test --home $CHAINDIR/$CHAINID2 --chain-id=$CHAINID2

echo "Adding genesis accounts..."
echo $MNEMONIC1 | $BINARY keys add val1 --home $CHAINDIR/$CHAINID1 --recover --keyring-backend=test 
echo $MNEMONIC2 | $BINARY keys add val2 --home $CHAINDIR/$CHAINID2 --recover --keyring-backend=test 
$BINARY add-genesis-account $($BINARY --home $CHAINDIR/$CHAINID1 keys show val1 --keyring-backend test -a) 100000000000stake  --home $CHAINDIR/$CHAINID1
$BINARY add-genesis-account $($BINARY --home $CHAINDIR/$CHAINID2 keys show val2 --keyring-backend test -a) 100000000000stake  --home $CHAINDIR/$CHAINID2

echo "Creating and collecting gentx..."
$BINARY gentx val1 7000000000stake --home $CHAINDIR/$CHAINID1 --chain-id $CHAINID1 --keyring-backend test
$BINARY gentx val2 7000000000stake --home $CHAINDIR/$CHAINID2 --chain-id $CHAINID2 --keyring-backend test
$BINARY collect-gentxs --home $CHAINDIR/$CHAINID1
$BINARY collect-gentxs --home $CHAINDIR/$CHAINID2

echo "Changing defaults and ports in app.toml and config.toml files..."
sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT1"'"#g' $CHAINDIR/$CHAINID1/config/config.toml
sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT1"'"#g' $CHAINDIR/$CHAINID1/config/config.toml
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAINDIR/$CHAINID1/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAINDIR/$CHAINID1/config/config.toml
sed -i '' 's/index_all_keys = false/index_all_keys = true/g' $CHAINDIR/$CHAINID1/config/config.toml
sed -i '' 's/enable = false/enable = true/g' $CHAINDIR/$CHAINID1/config/app.toml
sed -i '' 's/swagger = false/swagger = true/g' $CHAINDIR/$CHAINID1/config/app.toml
sed -i '' 's#"tcp://0.0.0.0:1317"#"tcp://0.0.0.0:'"$RESTPORT1"'"#g' $CHAINDIR/$CHAINID1/config/app.toml

sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT2"'"#g' $CHAINDIR/$CHAINID2/config/config.toml
sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT2"'"#g' $CHAINDIR/$CHAINID2/config/config.toml
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAINDIR/$CHAINID2/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAINDIR/$CHAINID2/config/config.toml
sed -i '' 's/index_all_keys = false/index_all_keys = true/g' $CHAINDIR/$CHAINID2/config/config.toml
sed -i '' 's/enable = false/enable = true/g' $CHAINDIR/$CHAINID2/config/app.toml
sed -i '' 's/swagger = false/swagger = true/g' $CHAINDIR/$CHAINID2/config/app.toml
sed -i '' 's#"tcp://0.0.0.0:1317"#"tcp://0.0.0.0:'"$RESTPORT2"'"#g' $CHAINDIR/$CHAINID2/config/app.toml

# Start chain 1
echo "Starting $CHAINID1 in $CHAINDIR..."
echo "Creating log file at $CHAINDIR/$CHAINID1.log"
$BINARY start --home $CHAINDIR/$CHAINID1 --pruning=nothing --grpc.address="0.0.0.0:$GRPCPORT1" > $CHAINDIR/$CHAINID1.log 2>&1 &

# Start chain 2
echo "Starting $CHAINID1 in $CHAINDIR..."
echo "Creating log file at $CHAINDIR/$CHAINID1.log"
$BINARY start --home $CHAINDIR/$CHAINID2 --pruning=nothing --grpc.address="0.0.0.0:$GRPCPORT2" > $CHAINDIR/$CHAINID2.log 2>&1 &


