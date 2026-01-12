SHELL := /bin/bash

BINARY := godoc
PKG := github.com/liasica/godoc
OUTDIR := bin

# Try to infer values from git/date, fall back to "unknown"/v1.0.0 when unavailable
TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo v1.0.0)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
# detect dirty working tree
DIRTY := $(shell test -n "$$(shell git status --porcelain 2>/dev/null)" && echo -dirty || echo)

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

# List of platforms to build for when running `make build`.
# You can customize this list or pass a single GOOS/GOARCH to `make build-local`.
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

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

.PHONY: all build build-local clean
all: build

# Default `build` will produce binaries for the platforms listed in PLATFORMS.
# Each output will be placed in $(OUTDIR) and named $(BINARY)-<os>-<arch>[.exe]
build:
	@echo "Building $(BINARY) for platforms: $(PLATFORMS)"
	mkdir -p $(OUTDIR)
	@for platform in $(PLATFORMS); do \
		os="$${platform%/*}"; arch="$${platform#*/}"; \
		ext=""; if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		out="$(OUTDIR)/$(BINARY)-$$os-$$arch$$ext"; \
		echo "  -> $$out"; \
		GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags "$(LDFLAGS)" -o $$out ./cmd/godoc || exit 1; \
	done

# Build only for the host (or for GOOS/GOARCH you set in the environment).
build-local:
	@echo "Building $(BINARY) for $(GOOS)/$(GOARCH)"
	mkdir -p $(OUTDIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -ldflags "$(LDFLAGS)" -o $(OUTDIR)/$(BINARY) ./cmd/godoc
	@echo "Built $(OUTDIR)/$(BINARY)"

clean:
	rm -rf $(OUTDIR)
