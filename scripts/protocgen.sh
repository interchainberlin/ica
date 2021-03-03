#!/usr/bin/env bash

set -eo pipefail
# GO=/home/sean/go/bin/go

# Get the path of the cosmos-sdk repo from go/pkg/mod
# cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
proto_dirs=$(find . -path ./third_party -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  buf protoc \
  -I "proto" \
  -I "third_party/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,\
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  $(find "${dir}" -name '*.proto')
done

cp -r github.com/interchainberlin/interchainaccount/* ./
rm -rf github.com
