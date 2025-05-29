# VERSION is based on a date stamp plus the last commit
VERSION?=v$(shell date +%Y%m%d)-$(shell git describe --tags)
BRANCH?=$(shell git branch --show-current)
SHA1?=$(shell git rev-parse HEAD)
BUILD=$(shell date +%FT%T%z)
LDFLAG_LOCATION=github.com/yshngg/kle/pkg/version
LDFLAGS=-ldflags "-X ${LDFLAG_LOCATION}.version=${VERSION} -X ${LDFLAG_LOCATION}.buildDate=${BUILD} -X ${LDFLAG_LOCATION}.gitbranch=${BRANCH} -X ${LDFLAG_LOCATION}.gitsha1=${SHA1}"
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# CURRENT_DIR is the current dir where the Makefile exists
CURRENT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: pwd
pwd:
	@echo $(CURRENT_DIR)

.PHONY: all
all: build

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build ${LDFLAGS} -o _output/bin/kle github.com/yshngg/kle/cmd

.PHONY: clean
clean:
	rm -rf _output
	rm -rf _tmp
