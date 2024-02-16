#!/bin/bash

set -e -x -o pipefail

install_forge_cli() {
  echo "Installing forge-cli..."
  # Determine system architecture
  arch=$(uname -m)
  case $arch in
  x86_64 | amd64)
    suffix=""
    ;;
  aarch64)
    suffix=-arm64
    ;;
  armv7l)
    suffix=-armhf
    ;;
  *)
    echo "Unsupported architecture $arch"
    exit 1
    ;;
  esac

  # Download the appropriate binary
  $SUDO curl -fsSL https://github.com/forge4flow/forge-cli/releases/latest/download/forge-cli${suffix} --output /usr/local/bin/forge-cli
  $SUDO chmod +x /usr/local/bin/forge-cli
}

install_forge_cli