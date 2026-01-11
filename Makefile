SHELL := /bin/bash

BINARY := godoc
PKG := github.com/liasica/godoc
OUTDIR := bin

# Try to infer values from git/date, fall back to "unknown"/v1.0.0 when unavailable
TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo v1.0.0)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
# detect dirty working tree
DIRTY := $(shell test -n "$$$(shell git status --porcelain 2>/dev/null)" && echo -dirty || echo)

# Allow overriding the base version (e.g. from CI). Default base is the latest tag.
VERSION_BASE ?= $(TAG)
VERSION := $(VERSION_BASE)-$(COMMIT)$(DIRTY)

LDFLAGS := -s -w \
	-X '$(PKG).Version=$(VERSION)' \
	-X '$(PKG).BuildTime=$(BUILDTIME)' \
	-X '$(PKG).CommitHash=$(COMMIT)'

# Allow overriding GOOS/GOARCH from environment, default to host
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

## Print all variables for debugging in a Makefile-safe way
## Use $(info ...) so make prints these during evaluation (no leading tabs or shell @ required)
$(info ----------------- Variables ----------------)
$(info TAG: $(TAG))
$(info COMMIT: $(COMMIT))
$(info BUILDTIME: $(BUILDTIME))
$(info DIRTY: $(DIRTY))
$(info VERSION_BASE: $(VERSION_BASE))
$(info VERSION: $(VERSION))
$(info LDFLAGS: $(LDFLAGS))
$(info GOOS: $(GOOS))
$(info GOARCH: $(GOARCH))

.PHONY: all build clean
all: build

build:
	@echo "Building $(BINARY) for $(GOOS)/$(GOARCH)"
	mkdir -p $(OUTDIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -ldflags "$(LDFLAGS)" -o $(OUTDIR)/$(BINARY) ./cmd/godoc
	@echo "Built $(OUTDIR)/$(BINARY)"

clean:
	rm -rf $(OUTDIR)
