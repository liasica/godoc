#!/usr/bin/env bash

set -e

# GitHub repository information
GITHUB_REPO="liasica/godoc"
BINARY_NAME="godoc"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to extract version number from godoc --version output
get_installed_version() {
    if command -v godoc >/dev/null 2>&1; then
        local version_output
        version_output=$(godoc --version 2>/dev/null || echo "")
        if [ -n "$version_output" ]; then
            # Extract version from "godoc version godoc has version v1.0.3-b2b0166"
            local version
            version=$(echo "$version_output" | grep "godoc version" | sed -n 's/.*version \(v[0-9.]*\).*/\1/p')
            echo "$version"
        else
            echo ""
        fi
    else
        echo ""
    fi
}

# Function to get the latest release version from GitHub
get_latest_version() {
    local latest_version
    latest_version=$(curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$latest_version" ]; then
        print_error "Failed to fetch the latest version from GitHub"
        exit 1
    fi
    echo "$latest_version"
}

# Function to compare versions (returns 0 if v1 < v2, 1 if v1 >= v2)
version_lt() {
    local v1=$1
    local v2=$2
    # Remove 'v' prefix
    v1=${v1#v}
    v2=${v2#v}

    # Use sort -V for version comparison
    if [ "$(printf '%s\n' "$v1" "$v2" | sort -V | head -n1)" = "$v1" ] && [ "$v1" != "$v2" ]; then
        return 0
    else
        return 1
    fi
}

# Function to detect OS and architecture
detect_platform() {
    local os
    local arch
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)

    case "$os" in
        linux*)
            os="linux"
            ;;
        darwin*)
            os="darwin"
            ;;
        mingw* | msys* | cygwin*)
            os="windows"
            ;;
        *)
            print_error "Unsupported OS: $os"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64 | amd64)
            arch="amd64"
            ;;
        arm64 | aarch64)
            arch="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac

    echo "${os}-${arch}"
}

# Function to download and install the binary
install_binary() {
    local version=$1
    local platform=$2

    # Construct download URL and binary name
    local binary_file="${BINARY_NAME}-${platform}"
    if [[ "$platform" == windows* ]]; then
        binary_file="${binary_file}.exe"
    fi

    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${binary_file}"

    print_info "Downloading ${BINARY_NAME} ${version} for ${platform}..."

    # Get GOPATH
    local gopath
    gopath=$(go env GOPATH)
    if [ -z "$gopath" ]; then
        print_error "GOPATH is not set. Please install Go first."
        exit 1
    fi

    local install_dir="${gopath}/bin"
    mkdir -p "$install_dir"

    local target_file="${install_dir}/${BINARY_NAME}"
    if [[ "$platform" == windows* ]]; then
        target_file="${target_file}.exe"
    fi

    # Download the binary
    if ! curl -fsSL -o "$target_file" "$download_url"; then
        print_error "Failed to download ${binary_file} from ${download_url}"
        exit 1
    fi

    # Make it executable
    chmod +x "$target_file"

    print_info "${BINARY_NAME} ${version} has been installed to ${target_file}"
}

# Main installation logic
main() {
    print_info "Checking ${BINARY_NAME} installation..."

    # Step 1: Check if godoc is already installed
    local installed_version
    installed_version=$(get_installed_version)
    if [ -n "$installed_version" ]; then
        print_info "Found installed version: ${installed_version}"
    else
        print_warn "${BINARY_NAME} is not installed or version cannot be determined"
    fi

    # Step 2: Get the latest release version from GitHub
    print_info "Fetching the latest version from GitHub..."
    local latest_version
    latest_version=$(get_latest_version)
    print_info "Latest version available: ${latest_version}"

    # Step 3: Compare versions and install if necessary
    if [ -z "$installed_version" ]; then
        print_info "Installing ${BINARY_NAME} ${latest_version}..."
        local platform
        platform=$(detect_platform)
        install_binary "$latest_version" "$platform"
        print_info "Installation complete!"
    elif version_lt "$installed_version" "$latest_version"; then
        print_info "Upgrading from ${installed_version} to ${latest_version}..."
        local platform
        platform=$(detect_platform)
        install_binary "$latest_version" "$platform"
        print_info "Upgrade complete!"
    else
        print_info "${BINARY_NAME} is already up to date (${installed_version})"
    fi

    # Verify installation
    print_info "Verifying installation..."
    if command -v godoc >/dev/null 2>&1; then
        godoc --version
    else
        print_warn "Please add \$(go env GOPATH)/bin to your PATH"
    fi
}

main

