include .env
export

PROJECT_NAME := NexusModsWatcher

GOBIN := ./bin

STDERR := /tmp/$(PROJECT_NAME)-stderr.txt
MAKEFLAGS += --silent

MAIN_FILENAME := src/main.go
BINARY_FILENAME := watcher

ifeq ($(RELEASE), true)
	BUILD_FLAGS := -ldflags "-X main.Release=true"
else
	# https://golang.org/cmd/compile/#hdr-Command_Line
	# https://golang.org/doc/gdb#Introduction
	BUILD_FLAGS := -gcflags=all="-N -l"
endif

.PHONY: help
help: Makefile
	@echo "Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n "s/^##//p" $< | column -t -s ":" |  sed -e "s/^/ /"
	@echo

#  Make Commands
# --------------------------------------------------
## build: build project
.PHONY: build
build: go-build

## clean: clean up caches
.PHONY: clean
clean: go-clean

## compile: clean up caches, resolve dependencies, and build the application
.PHONY: compile
compile:
	@-rm -f $(STDERR)
	@-touch $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/^/Error:\n/' | sed '/^make\[.*\]/d' | sed "s/^/  /g" | sed -r "s/(.*)/\x1b[38;5;9m\1\x1b[0m/g"

## fmt: run formatter
.PHONY: fmt
fmt: go-fmt

## gen: generate source code
.PHONY: gen
gen: go-generate

## lint: run linters (golangci-lint)
.PHONY: lint
lint: go-lint

## resolve: resolve dependencies
.PHONY: resolve
resolve: go-mod

## test: run all tests
.PHONY: test
test:
	@go test -race -covermode=atomic -coverprofile=coverage.out \
		$(dir $(abspath $(firstword $(MAKEFILE_LIST))))src/log \
		$(dir $(abspath $(firstword $(MAKEFILE_LIST))))src/nexus

#  Go Commands
# --------------------------------------------------
.PHONY: go-build
go-build:
	@echo "▶ building binary: BUILD_FLAGS = $(BUILD_FLAGS)"
	@mkdir -p $(GOBIN)
	@go build $(BUILD_FLAGS) -o $(GOBIN)/$(BINARY_FILENAME) $(MAIN_FILENAME)

.PHONY: go-clean
go-clean:
	@echo "▶ cleaning up caches"
	@go clean -cache -testcache

.PHONY: go-compile
go-compile: go-mod go-vet go-build

.PHONY: go-fmt
go-fmt:
	@echo "▶ reformatting code"
	@go fmt $(dir $(abspath $(firstword $(MAKEFILE_LIST))))src/...

.PHONY: go-generate
go-generate:
	@echo "▶ generating code"
	@go generate $(dir $(abspath $(firstword $(MAKEFILE_LIST))))src/...

.PHONY: go-lint
go-lint:
	@echo "▶ running linters"
	@golangci-lint run

.PHONY: go-mod
go-mod:
	@echo "▶ checking if there is any missing dependencies"
	@go mod tidy

.PHONY: go-vet
go-vet:
	@echo "▶ vetting source code"
	@go vet $(dir $(abspath $(firstword $(MAKEFILE_LIST))))src/...
