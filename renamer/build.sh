#!/bin/sh
# Builds renamer for Windows, macOS and Linux into ./dist
set -e
cd "$(dirname "$0")"
mkdir -p dist
GOOS=windows GOARCH=amd64 go build -o dist/renamer-windows.exe .
GOOS=darwin  GOARCH=arm64 go build -o dist/renamer-mac-apple-silicon .
GOOS=darwin  GOARCH=amd64 go build -o dist/renamer-mac-intel .
GOOS=linux   GOARCH=amd64 go build -o dist/renamer-linux .
echo "Built into dist/"
