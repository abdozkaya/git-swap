#!/bin/bash

# Configuration
REPO="abdozkaya/git-swap"
BIN_NAME="git-swap"
INSTALL_DIR="/usr/local/bin"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "🔎 Detecting system..."

# Detect OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS_TYPE="linux";;
    Darwin*)    OS_TYPE="darwin";;
    *)          echo -e "${RED}Unsupported OS: ${OS}${NC}"; exit 1;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)    ARCH_TYPE="amd64";;
    arm64)     ARCH_TYPE="arm64";;
    aarch64)   ARCH_TYPE="arm64";;
    *)         echo -e "${RED}Unsupported architecture: ${ARCH}${NC}"; exit 1;;
esac

# Construct Download URL based on release naming convention
TARGET_FILE="${BIN_NAME}-${OS_TYPE}-${ARCH_TYPE}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${TARGET_FILE}"

echo "⬇️  Downloading ${BIN_NAME} for ${OS_TYPE}/${ARCH_TYPE}..."
curl -L -o ${BIN_NAME} ${DOWNLOAD_URL} --fail

if [ $? -ne 0 ]; then
    echo -e "${RED}Download failed! Please ensure the release exists on GitHub.${NC}"
    exit 1
fi

chmod +x ${BIN_NAME}

echo "📦 Installing to ${INSTALL_DIR}..."
if [ -w "${INSTALL_DIR}" ]; then
    mv ${BIN_NAME} ${INSTALL_DIR}/${BIN_NAME}
else
    sudo mv ${BIN_NAME} ${INSTALL_DIR}/${BIN_NAME}
fi

echo -e "${GREEN}✅ git-swap installed successfully!${NC}"
echo "Run 'git-swap help' to get started."