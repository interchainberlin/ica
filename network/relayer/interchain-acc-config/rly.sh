#!/bin/bash
DEV_MNEMONIC_1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
DEV_MNEMONIC_2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"

rly config init

rly config add-chains $PWD/network/relayer/interchain-acc-config/chains
rly config add-paths $PWD/network/relayer/interchain-acc-config/paths

rly keys restore test-1 test-1 "$DEV_MNEMONIC_1"

rly keys restore test-2 test-2 "$DEV_MNEMONIC_2"

rly light init test-1 -f
rly light init test-2 -f

rly tx path test1-account-test2

rly start test1-account-test2
