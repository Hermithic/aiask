#!/bin/bash
# AIask Installer Script
# Usage: curl -fsSL https://raw.githubusercontent.com/Hermithic/aiask/master/install.sh | bash

set -e

VERSION="2.0.1"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="aiask"

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    linux|darwin)
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

DOWNLOAD_URL="https://github.com/Hermithic/aiask/releases/download/v${VERSION}/aiask-${VERSION}-${OS}-${ARCH}.tar.gz"

echo "🤖 Installing AIask v${VERSION}..."
echo "   OS: ${OS}, Arch: ${ARCH}"
echo ""

# Create temp directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download
echo "📥 Downloading from ${DOWNLOAD_URL}..."
curl -fsSL "$DOWNLOAD_URL" -o aiask.tar.gz

# Extract
echo "📦 Extracting..."
tar -xzf aiask.tar.gz

# Install
echo "🔧 Installing to ${INSTALL_DIR}..."
if [ -w "$INSTALL_DIR" ]; then
    mv "aiask-${OS}-${ARCH}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    sudo mv "aiask-${OS}-${ARCH}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

# Cleanup
cd /
rm -rf "$TMP_DIR"

# Verify
if command -v aiask &> /dev/null; then
    echo ""
    echo "✅ AIask installed successfully!"
    echo ""
    aiask version
    echo ""
    echo "Run 'aiask config' to set up your AI provider."
else
    echo "❌ Installation failed. Please check your PATH."
    exit 1
fi

