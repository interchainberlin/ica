# Interchain Accounts
## Local Demo

### Setup

```bash
# Clone this repository and build
git clone https://github.com/cosmos/interchain-accounts.git
cd interchain-accounts
make install 

# Hermes Relayer
[Hermes](https://hermes.informal.systems/) is a Rust implementation of a relayer for the [Inter-Blockchain Communication (IBC)](https://ibcprotocol.org/) protocol.

- Install [Rust](https://www.rust-lang.org/tools/install)
- Install [Hermes](https://hermes.informal.systems/installation.html)

# Bootstrap two local chains & create a connection using the hermes relayer
make init

# Wait for the ibc connection & channel handshake to complete and the relayer to start
```

### Demo

```bash
# Open a seperate terminal

# These address are defined in init.sh for development purposes
export VAL_1=cosmos1mjk79fjjgpplak5wq838w0yd982gzkyfrk07am
export VAL_2=cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs

# Register an IBC Account on chain test-2 
icad tx intertx register --from val1 --connection-id connection-0 --chain-id test-1 --gas 150000 --home ./data/test-1 --node tcp://localhost:16657 --keyring-backend test -y

# Start the hermes relayer in the first terminal
# This will also finish the channel creation handshake signalled during the register step
make start-rly

# Get the address of interchain account
icad query ibcaccount address cosmos1mjk79fjjgpplak5wq838w0yd982gzkyfrk07am connection-0 --home ./data/test-2 --node tcp://localhost:26657
# Output -> account_address: cosmos1plyxrjdepap2zgqmfpzfchmklwqhchq5jrctm0

export IBC_ACCOUNT=cosmos1plyxrjdepap2zgqmfpzfchmklwqhchq5jrctm0

# Check the interchain account's balance on test-2 chain. It should be empty.
icad q bank balances $IBC_ACCOUNT --chain-id test-2 --node tcp://localhost:26657

# Send some assets to $IBC_ACCOUNT.
icad tx bank send val2 $IBC_ACCOUNT 1000stake --chain-id test-2 --home ./data/test-2 --node tcp://localhost:26657 --keyring-backend test -y

# Check that the balance has been updated
icad q bank balances $IBC_ACCOUNT --chain-id test-2 --node tcp://localhost:26657

# Test sending assets from interchain account via ibc.
icad tx intertx send cosmos1plyxrjdepap2zgqmfpzfchmklwqhchq5jrctm0 $VAL_2 500stake --connection-id conection-0 --chain-id test-1 --gas 90000 --home ./data/test-1 --node tcp://localhost:16657 --from val1 --keyring-backend test -y

# Wait until the relayer has relayed the packet

# Check if the balance has been changed (it should now be 500stake)
icad q bank balances $IBC_ACCOUNT --chain-id test-2 --node tcp://localhost:26657
```

## Collaboration

Please use conventional commits  https://www.conventionalcommits.org/en/v1.0.0/

```
chore(bump): bumping version to 2.0
fix(bug): fixing issue with...
feat(featurex): adding feature...
```
