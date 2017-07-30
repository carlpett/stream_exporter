GO           := go
FIRST_GOPATH := $(firstword $(subst :, ,$(GOPATH)))
PROMU        := $(FIRST_GOPATH)/bin/promu
pkgs          = $(shell $(GO) list ./... | grep -v /vendor/)

PREFIX                  ?= $(shell pwd)/.build
BIN_DIR                 ?= $(shell pwd)

ifdef DEBUG
	bindata_flags = -debug
endif

all: format build

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

test:
	@echo ">> testing code"
	@$(GO) test $(pkgs)

build: promu
	@echo ">> building binaries"
	@$(PROMU) build --prefix $(PREFIX)

tarball: promu
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)

promu:
	@GOOS=$(shell uname -s | tr A-Z a-z) \
	GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m))) \
	$(GO) get -u github.com/prometheus/promu


.PHONY: all style format build vet tarball docker promu