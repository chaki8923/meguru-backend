#!/bin/bash

set -e

echo "Building Gemini Lambda function..."

# Set environment for Linux build
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# Build the binary
go mod tidy
go build -ldflags="-s -w" -o bootstrap main.go

# Create ZIP file
zip -j gemini-lambda.zip bootstrap

# Move ZIP to terraform directory
mv gemini-lambda.zip ../../terraform/

echo "Lambda function built successfully: ../../terraform/gemini-lambda.zip"

# Clean up
rm bootstrap

echo "Build completed!"
