PREFIX ?= /usr/local
BIN ?= $(PREFIX)/bin
UNAME_S ?= $(shell uname -s)
UNAME_M ?= $(shell uname -m)

BUF_VERSION ?= 0.11.0

PROTOC_VERSION ?= 3.11.2
ifeq ($(UNAME_S),Linux)
  PROTOC_ZIP ?= protoc-3.11.2-linux-x86_64.zip
endif
ifeq ($(UNAME_S),Darwin)
  PROTOC_ZIP ?= protoc-3.11.2-osx-x86_64.zip
endif

GO=/home/sean/go/bin/go

proto-tools: proto-tools-stamp
proto-tools-stamp:
	@echo "Installing protoc compiler..."
	@(cd /tmp; \
	curl -OL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"; \
	unzip -o ${PROTOC_ZIP} -d $(PREFIX) bin/protoc; \
	unzip -o ${PROTOC_ZIP} -d $(PREFIX) 'include/*'; \
	rm -f ${PROTOC_ZIP})

	@echo "Installing protoc-gen-buf-check-breaking..."
	@curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/protoc-gen-buf-check-breaking-${UNAME_S}-${UNAME_M}" \
    -o "${BIN}/protoc-gen-buf-check-breaking" && \
	chmod +x "${BIN}/protoc-gen-buf-check-breaking"

	@echo "Installing buf..."
	@curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-${UNAME_S}-${UNAME_M}" \
    -o "${BIN}/buf" && \
	chmod +x "${BIN}/buf"

	touch $@

protoc-gen-gocosmos:
	@echo "Installing protoc-gen-gocosmos..."
	@ $(GO) install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos

proto-gen:
	@./scripts/protocgen.sh

install: go.sum
	go install ./demo/cmd/demod
	go install ./demo/cmd/democli