GO := go
PROMU := $(GOPATH)/bin/promu
pkgs = $(shell go list ./... | grep -v /vendor/)

all: format build

build:
	$(PROMU) build

format:
	$(GO) fmt $(pkgs)

.PHONY: all build format
