#!/usr/bin/env bash

set -e

# GitKit Installation Script
# This script automatically detects your OS and installs GitKit

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Version to install (latest if not specified)
VERSION="${GITKIT_VERSION:-latest}"
GITHUB_REPO="samanar/gitkit"

echo -e "${GREEN}GitKit Installation Script${NC}"
echo "========================================"

# Detect OS and Architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        mingw* | msys* | cygwin*)
            OS="windows"
            ;;
        *)
            echo -e "${RED}Unsupported operating system: $OS${NC}"
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64 | amd64)
            ARCH="amd64"
            ;;
        aarch64 | arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac

    echo -e "Detected OS: ${GREEN}$OS${NC}"
    echo -e "Detected Architecture: ${GREEN}$ARCH${NC}"
}

# Get latest version from GitHub
get_latest_version() {
    if [ "$VERSION" = "latest" ]; then
        echo -e "${YELLOW}Fetching latest version...${NC}"
        VERSION=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [ -z "$VERSION" ]; then
            echo -e "${RED}Failed to fetch latest version${NC}"
            exit 1
        fi
    fi
    echo -e "Version to install: ${GREEN}$VERSION${NC}"
}

# Download binary
download_binary() {
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="gitkit-${OS}-${ARCH}.exe"
    else
        BINARY_NAME="gitkit-${OS}-${ARCH}"
    fi

    DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$VERSION/$BINARY_NAME"
    
    echo -e "${YELLOW}Downloading GitKit...${NC}"
    echo "URL: $DOWNLOAD_URL"
    
    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    TMP_FILE="$TMP_DIR/gitkit"
    
    if command -v curl &> /dev/null; then
        curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"
    elif command -v wget &> /dev/null; then
        wget -O "$TMP_FILE" "$DOWNLOAD_URL"
    else
        echo -e "${RED}Neither curl nor wget found. Please install one of them.${NC}"
        exit 1
    fi
    
    if [ ! -f "$TMP_FILE" ]; then
        echo -e "${RED}Failed to download binary${NC}"
        exit 1
    fi
    
    chmod +x "$TMP_FILE"
    echo -e "${GREEN}Download completed${NC}"
}

# Install binary
install_binary() {
    echo -e "${YELLOW}Installing GitKit...${NC}"
    
    # Determine installation directory
    if [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    elif [ -w "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    elif [ -w "$HOME/bin" ]; then
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
    else
        # Need sudo for system-wide installation
        INSTALL_DIR="/usr/local/bin"
        USE_SUDO=true
    fi
    
    INSTALL_PATH="$INSTALL_DIR/gitkit"
    
    # Install with or without sudo
    if [ "$USE_SUDO" = true ]; then
        echo -e "${YELLOW}Installation requires sudo privileges...${NC}"
        sudo mv "$TMP_FILE" "$INSTALL_PATH"
        sudo chmod +x "$INSTALL_PATH"
    else
        mv "$TMP_FILE" "$INSTALL_PATH"
        chmod +x "$INSTALL_PATH"
    fi
    
    # Clean up
    rm -rf "$TMP_DIR"
    
    echo -e "${GREEN}✓ GitKit installed successfully to: $INSTALL_PATH${NC}"
}

# Verify installation
verify_installation() {
    echo -e "${YELLOW}Verifying installation...${NC}"
    
    if command -v gitkit &> /dev/null; then
        INSTALLED_VERSION=$(gitkit version 2>/dev/null || echo "installed")
        echo -e "${GREEN}✓ GitKit is ready to use!${NC}"
        echo ""
        echo "Run 'gitkit --help' to get started"
    else
        echo -e "${YELLOW}Warning: gitkit command not found in PATH${NC}"
        echo "You may need to add $INSTALL_DIR to your PATH"
        echo ""
        echo "Add this line to your shell configuration file (~/.bashrc, ~/.zshrc, etc.):"
        echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    fi
}

# Main installation flow
main() {
    detect_platform
    get_latest_version
    download_binary
    install_binary
    verify_installation
    
    echo ""
    echo "========================================"
    echo -e "${GREEN}Installation complete!${NC}"
    echo ""
    echo "Quick start:"
    echo "  gitkit init              # Initialize GitKit in your repo"
    echo "  gitkit feature start     # Start a new feature branch"
    echo "  gitkit --help            # Show all commands"
    echo ""
}

# Run main function
main
