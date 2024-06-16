#!/bin/bash

REPO_URL="https://github.com/ReallyLiri/kubectui.git"
CLONE_DIR=$(mktemp -d)

cleanup() {
  echo "Cleaning up..."
  rm -rf "$CLONE_DIR"
}

trap cleanup EXIT

git clone "$REPO_URL" "$CLONE_DIR"

cd "$CLONE_DIR"
go install .
echo "Installation complete."
