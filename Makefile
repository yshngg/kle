# Makefile for kle project

# Compute metadata
VERSION ?= $(shell git describe --tags --always)-dev-$(shell git rev-parse --short HEAD)
BRANCH  ?= $(shell git rev-parse --abbrev-ref HEAD)
SHA1    ?= $(shell git rev-parse --short HEAD)
BUILD   ?= $(shell date -u +%FT%T%z)

# Linker flags
LDFLAG_LOCATION = github.com/yshngg/kle/pkg/version
LDFLAGS = -ldflags "-s -w \
	-X ${LDFLAG_LOCATION}.version=${VERSION} \
	-X ${LDFLAG_LOCATION}.buildDate=${BUILD} \
	-X ${LDFLAG_LOCATION}.gitbranch=${BRANCH} \
	-X ${LDFLAG_LOCATION}.gitsha1=${SHA1}"

# Defaults from go env
GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Directories
BIN_DIR   := _output/bin
CMD_PATH  := github.com/yshngg/kle/cmd
BINARY    := kle

# Helpers
define info
	@printf "INFO %s\n" "$(1)"
endef

.PHONY: all build clean fmt vet lint test
all: fmt vet lint test build

build: $(BIN_DIR)/$(BINARY)

$(BIN_DIR)/$(BINARY):
	$(call info,Building ${BINARY} ${VERSION} for ${GOOS}/${GOARCH}...)
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build ${LDFLAGS} -o $@ ${CMD_PATH}

clean:
	$(call info,Cleaning build artifacts...)
	rm -rf _output
	@echo Done.

fmt:
	$(call info,Formatting code...)
	go fmt ./...

vet:
	$(call info,Running go vet...)
	go vet ./...

lint:
	$(call info,Linting code with golangci-lint...)
	golangci-lint run

test:
	$(call info,Running tests...)
	go test ./... -timeout 30s
