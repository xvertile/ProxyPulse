#!/bin/bash
# Build for the current platform
# Set environment variables for cross-compilation
export GOARCH=amd64
export GOOS=linux

# Build for Linux platform and output to build directory
go build -o build/proxypulse

# Reset environment variables (optional)
unset GOARCH
unset GOOS
