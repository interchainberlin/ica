#!/bin/bash
set -e

### Sleep is needed otherwise the relayer crashes when trying to init
sleep 1s
### Restore Keys
hermes -c ./network/hermes/config.toml keys restore test-2 -m "record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"
sleep 5s
hermes -c ./network/hermes/config.toml keys restore test-1 -m "alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
sleep 5s

