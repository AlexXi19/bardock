#!/bin/bash

set -e

REPO="alexxi19/bardock"
LATEST_RELEASE=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep "tag_name" | awk -F '"' '{print $4}')
DOWNLOAD_FILE="bardock.zip"
EXECUTABLE_NAME="bardock"

OS=$(uname)

if [ "$OS" == "Darwin" ]; then
  DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/bardock-$LATEST_RELEASE-darwin-arm64.tar.gz"
  echo "Detected macOS. Downloading latest release: $LATEST_RELEASE..."
elif [ "$OS" == "Linux" ]; then
  DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/bardock-$LATEST_RELEASE-linux-arm64.tar.gz"
  echo "Detected Linux. Downloading latest release: $LATEST_RELEASE..."
else
  echo "Unsupported OS. Exiting."
  exit 1
fi

# Download latest release
curl -L -o $DOWNLOAD_FILE $DOWNLOAD_URL

# Unzip the downloaded file
echo "Unzipping downloaded file..."
tar -xzf $DOWNLOAD_FILE

# Make it executable
chmod +x bardock

# Move it to /usr/local/bin or any other directory in $PATH
sudo mv $EXECUTABLE_NAME /usr/local/bin

echo "$EXECUTABLE_NAME installed successfully."
