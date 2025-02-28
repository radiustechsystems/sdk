#!/bin/bash
set -e

echo "Running cleanup for all Radius SDKs..."

echo -e "\nRunning Go SDK cleanup..."
cd go
./cleanup.sh
cd ..

echo -e "\nRunning TypeScript SDK cleanup..."
cd typescript
./cleanup.sh
cd ..

echo -e "\nAll SDK checks completed successfully!"
