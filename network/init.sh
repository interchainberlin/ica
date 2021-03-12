#!/bin/bash
DEV_MNEMONIC_1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
DEV_MNEMONIC_2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"

### Chain 1
yes $DEV_MNEMONIC_1 | icad keys add val --recover --home ~/.demo-test-1

icad init --chain-id test-1 test-1 --home ~/.demo-test-1

sed -i -e 's/2665/1665/g' ~/.demo-test-1/config/config.toml
sed -i -e 's#localhost:6060#localhost:6061#g' ~/.demo-test-1/config/config.toml
sed -i -e 's#address = "0.0.0.0:9090"#address = "0.0.0.0:9092"#g' ~/.demo-test-1/config/app.toml

icad add-genesis-account $(icad keys show val -a --home ~/.demo-test-1) 100000000000stake --home ~/.demo-test-1
icad gentx val 7000000000stake --chain-id test-1 --home ~/.demo-test-1

icad collect-gentxs --home ~/.demo-test-1

### Chain 2
yes $DEV_MNEMONIC_2 | icad keys add val2 --recover --home ~/.demo-test-2

icad init --chain-id test-2 test-2 --home ~/.demo-test-2

icad add-genesis-account $(icad keys show val2 -a --home ~/.demo-test-2) 100000000000stake --home ~/.demo-test-2
icad gentx val2 7000000000stake --chain-id test-2 --home ~/.demo-test-2

icad collect-gentxs --home ~/.demo-test-2

#echo "Running both chains"
icad start --pruning nothing --home ~/.demo-test-1 &
icad start --pruning nothing --home ~/.demo-test-2 &
