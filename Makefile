BUILD_DIR ?= $(CURDIR)/bin

.PHONY: all build clean run

all: lint build

run:
	@go run ./test/main.go

BUILD_TARGETS := build install

build: BUIDL_ARGS=-o $(BUILD_DIR)/

$(BUILD_TARGETS): go.sum $(BUILD_DIR)/
	go $@ -mod=readonly $(BUIDL_ARGS) ./test/...

$(BUILD_DIR)/:
	mkdir -p $(BUILD_DIR)

go.sum: go.mod
	go mod verify
	go mod tidy

lint:
	golangci-lint run --out-format=tab ./...

loc:
	tokei .