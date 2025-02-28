#!/bin/bash
set -e

echo "Running linting, formatting, and type checking for TypeScript SDK..."
echo "Running linting with fixes..."
pnpm lint:fix
echo "Running formatting with fixes..."
pnpm format:fix
echo "Running type checking..."
pnpm tsc --noEmit
echo "TypeScript SDK checks completed successfully!"
