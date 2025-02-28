#!/bin/bash
set -e

echo "Running linting and formatting for Go SDK..."
echo "Running go fmt..."
go fmt ./...
echo "Running golangci-lint..."
golangci-lint run
echo "Go SDK checks completed successfully!"
