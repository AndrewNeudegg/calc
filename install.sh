#!/bin/sh
set -e

# calc installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/AndrewNeudegg/calc/main/install.sh | sh

REPO="AndrewNeudegg/calc"
BINARY_NAME="calc"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Get latest release
echo "Fetching latest release..."
RELEASE_URL="https://api.github.com/repos/$REPO/releases/latest"
ASSET_NAME="${BINARY_NAME}_${OS}_${ARCH}"

# Try to get download URL
DOWNLOAD_URL=$(curl -sL "$RELEASE_URL" | grep "browser_download_url.*$ASSET_NAME" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Could not find release for $OS/$ARCH"
    echo "Please build from source: go build -o calc ./cmd/calc"
    exit 1
fi

# Download and install
echo "Downloading $BINARY_NAME..."
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

curl -sL "$DOWNLOAD_URL" -o "$BINARY_NAME"
chmod +x "$BINARY_NAME"

# Install to /usr/local/bin if possible, otherwise to ~/bin
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
    echo "Installing to $INSTALL_DIR (add to PATH if needed)"
else
    echo "Installing to $INSTALL_DIR"
fi

mv "$BINARY_NAME" "$INSTALL_DIR/"
cd - > /dev/null
rm -rf "$TMP_DIR"

echo ""
echo "✓ $BINARY_NAME installed successfully!"
echo ""
echo "Try it out:"
echo "  $BINARY_NAME"
echo "  $BINARY_NAME -c \"12 gbp in dollars\""
