#!/usr/bin/make -f
PACKAGES=$(shell go list ./...)
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

###############################################################################
###                           Install                                       ###
###############################################################################
install: go.sum
		@echo "--> Installing icad"
		@go install ./cmd/icad

install-debug: go.sum
	go build -gcflags="all=-N -l" ./cmd/icad

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

test:
	@go test -mod=readonly $(PACKAGES) -cover

lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

###############################################################################
###                                Protobuf                                 ###
###############################################################################

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh

proto-lint:
	@$(DOCKER_BUF) check lint --error-format=json

###############################################################################
###                                Initialize                               ###
###############################################################################
start-dev: install
	./network/init.sh

start-rly:
	./network/relayer/interchain-acc-config/rly.sh

kill-dev:
	rm -rf ~/.demo-test-1
	rm -rf ~/.demo-test-2
	rm -rf ~/.relayer
	killall icad
