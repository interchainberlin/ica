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

# Bootstrap two local chains
make start-dev

# Use relayer to link the chains
make start-rly
```

### Demo

```bash
# Register an IBC Account on chain test-2 
icad tx intertx register ibcaccount channel-0 --from val --chain-id test-1 --gas 90000 --home ~/.demo-test-1 --node tcp://localhost:16657

# Query for an IBC Account registered on behalf of an address (returns the address of the ibc account registered on chain test-2)
icad query intertx ibcaccount <address used to register account (val is used above)> ibcaccount channel-0 --node tcp://localhost:16657

# Query for an IBC Account directly on chain test-2 by address (return value of the previous query)
icad query ibcaccount ibcaccount <address (output of previous command)> --node tcp://localhost:26657
```


## Collaboration

Please use conventional commits  https://www.conventionalcommits.org/en/v1.0.0/

```
chore(bump): bumping version to 2.0
fix(bug): fixing issue with...
feat(featurex): adding feature...
```

