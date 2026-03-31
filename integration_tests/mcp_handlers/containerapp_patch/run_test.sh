#!/bin/bash

# ContainerApp Patch Integration Test Runner
# This script runs the containerapp patch integration tests

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

# Change to the project root directory
cd "$PROJECT_ROOT"

echo "=========================================="
echo "ContainerApp Patch Integration Tests"
echo "=========================================="
echo "Project root: $PROJECT_ROOT"
echo "Test file: $SCRIPT_DIR/containerapp_patch_test.go"
echo ""

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "Error: .env file not found in project root"
    echo "Please create a .env file with the required environment variables"
    exit 1
fi

echo "Running integration tests..."
echo ""

# Run the test
go run "$SCRIPT_DIR/containerapp_patch_test.go"

echo ""
echo "=========================================="
echo "Tests completed successfully!"
echo "=========================================="
