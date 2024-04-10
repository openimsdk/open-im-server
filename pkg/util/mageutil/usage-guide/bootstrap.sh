#!/bin/bash

if [[ ":$PATH:" == *":$HOME/.local/bin:"* ]]; then
    TARGET_DIR="$HOME/.local/bin"
else
    TARGET_DIR="/usr/local/bin"
    echo "Using /usr/local/bin as the installation directory. Might require sudo permissions."
fi

if ! command -v mage &> /dev/null; then
    echo "Installing Mage to $TARGET_DIR ..."
    GOBIN=$TARGET_DIR go install github.com/magefile/mage@latest
fi

if ! command -v mage &> /dev/null; then
    echo "Mage installation failed."
    echo "Please ensure that $TARGET_DIR is in your \$PATH."
    exit 1
fi

echo "Mage installed successfully."

go mod download
