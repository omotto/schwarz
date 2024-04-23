ROOT_DIR              ?= $(shell git rev-parse --show-toplevel)
GO                    ?= go
GO_FLAGS              ?=
GO_BIN                ?= ./bin/main
GOTEST_ARGS           ?= -timeout=5m -cover
GOLANG_CI_ARGS        ?= --allow-parallel-runners --timeout=5m
BUF                   ?= buf
mainpath              ?= main
pkgs                  ?= $(shell $(GO) list ./...)

ALL_PROTOS ?= $(shell find $(ROOT_DIR) \
	-type f -path $(ROOT_DIR)/templates -prune \
	-o -type f -path $(ROOT_DIR)/proto -prune \
	-o -type f -iname '*.proto')

ifneq ("$(wildcard $(ROOT_DIR)/.golangci.yaml)","")
	GOLANG_CI_ARGS += --config $(ROOT_DIR)/.golangci.yaml
endif

ifneq ("$(wildcard $(ROOT_DIR)/buf.yaml)","")
	BUF_ARGS += --config $(ROOT_DIR)/buf.yaml
endif

setup:
	go install github.com/bufbuild/buf/cmd/buf@v1.29.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

setup-local: setup
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.1

compile: go-dependencies
ifneq ($(wildcard ./$(mainpath)),)
	$(GO) build $(GO_FLAGS) -o $(GO_BIN) ./$(mainpath)
else
	$(GO) build $(GO_FLAGS) ./...
endif

all: lint test build

build: compile

test: go-dependencies
	@echo ">> running tests"
	$(GO) test $(GOTEST_ARGS) $(pkgs)

go-mod-download: $(ROOT_DIR)/.go-mod-download.sentinel

$(ROOT_DIR)/.go-mod-download.sentinel: $(ROOT_DIR)/go.mod $(ROOT_DIR)/go.sum
	$(GO) mod download && touch $(ROOT_DIR)/.go-mod-download.sentinel

go-dependencies: go-mod-download generate-all-rpcs

generate-all-rpcs: $(patsubst %.proto,%.pb.go,$(ALL_PROTOS))

%.pb.go: %.proto
	@echo ">> generating $@"
	@cd $(ROOT_DIR) && $(BUF) generate && touch $@

format-all-protos:
	@echo ">> formatting all *.proto files"
	@cd $(ROOT_DIR) && $(BUF) format -w

check-format-all-protos:
	@echo ">> formatting all *.proto files"
	@cd $(ROOT_DIR) && $(BUF) format -d --exit-code

lint-all-protos:
	@echo ">> linting all *.proto files"
	@cd $(ROOT_DIR) && $(BUF) $(BUF_ARGS) lint

format: go-dependencies
	@echo ">> formatting code"
	@cd $(ROOT_DIR) && golangci-lint run $(GOLANG_CI_ARGS) --disable-all --enable goimports,gofumpt --allow-parallel-runners=false --fix $(abspath .)/...

check-format: go-dependencies
	@echo ">> formatting code"
	@cd $(ROOT_DIR) && golangci-lint run $(GOLANG_CI_ARGS) --disable-all --enable goimports,gofumpt --allow-parallel-runners=false $(abspath .)/...

lint-golangci-lint: go-dependencies
	@echo ">> linting"
	@cd $(ROOT_DIR) && golangci-lint run $(GOLANG_CI_ARGS) $(abspath .)/...

lint-proto: lint-all-protos check-format-all-protos

fix: format-all-protos format

lint: lint-proto lint-golangci-lint check-format

.PHONY: all generate-rpcs clean run go-mod-download generate-all-rpcs go-dependencies
