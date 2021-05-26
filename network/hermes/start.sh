#!/bin/bash

# Start the hermes relayer in multi-paths mode
CONFIG_DIR=./network/hermes/config.toml

echo "Starting hermes relayer..."
hermes -c $CONFIG_DIR start-multi

