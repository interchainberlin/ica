# Interchain Accounts
This repo contains an ongoing refactor/update of https://github.com/chainapsis/cosmos-sdk-interchain-account which is based on the [ics-27 spec.](https://github.com/cosmos/ics/tree/master/spec/ics-027-interchain-accounts)

## Local Demo

### Setup

```bash
# Clone this repository and build
git clone https://github.com/interchainberlin/ica.git
cd ica
make install 

# Clone relayer and build
https://github.com/cosmos/relayer.git
cd relayer
make install

# Bootstrap two local chains & start the relayer with development data
make init
```

### Send Asset Demo

```bash
# These address are defined in init.sh for development purposes
VAL_1=cosmos1mjk79fjjgpplak5wq838w0yd982gzkyfrk07am
VAL_2=cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs

# Register an IBC Account on chain test-2 
icad tx intertx register --from val1 --source-port ibcaccount --source-channel channel-0 --chain-id test-1 --gas 90000 --home ./data/test-1 --node tcp://localhost:16657

# Get the address of interchain account
icad query intertx ibcaccount $VAL_1 ibcaccount channel-0 --node tcp://localhost:16657
# Output -> address: cosmos1pt6ar8lawmvvq5haxc3l3zhjfl04u56fs2ndh9

IBC_ACCOUNT=cosmos1pt6ar8lawmvvq5haxc3l3zhjfl04u56fs2ndh9

# Check the interchain account's balance on test-2 chain. It should be empty.
icad q bank balances $IBC_ACCOUNT --chain-id test-2

# Send some assets to $IBC_ACCOUNT.
icad tx bank send val2 $IBC_ACCOUNT 1000stake --chain-id test-2 --home ./data/test-2 --node tcp://localhost:26657

# Test sending assets on intetchain account via ibc.
icad tx intertx send cosmos-sdk $VAL_2 500stake --source-port ibcaccount --source-channel channel-0 --chain-id test-1 --gas 90000 --home ./data/test-1 --node tcp://localhost:16657 --from val1

# Wait until the relayer has relayed the packet

# Check if the balance has been changed.
icad q bank balances $IBC_ACCOUNT --chain-id test-2
icad q bank balances $VAL_2 --chain-id test-2
```


## Collaboration

Please use conventional commits  https://www.conventionalcommits.org/en/v1.0.0/

```
chore(bump): bumping version to 2.0
fix(bug): fixing issue with...
feat(featurex): adding feature...
```


