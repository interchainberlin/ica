#!/bin/bash

rly config init

rly config add-chains ./chains
rly config add-paths ./paths

rly keys restore test-1 test-1 "<mnemonic>"

rly keys restore test-2 test-2 "<mnemonic>"

rly light init test-1 -f
rly light init test-2 -f

rly tx path test1-account-test2

rly start test1-account-test2

