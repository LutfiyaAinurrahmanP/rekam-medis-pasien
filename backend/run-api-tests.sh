#!/usr/bin/env bash

# Test runner script with Jest-like output
# This script runs all API tests and provides a nice summary

set -e

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║                  API Test Suite Runner                        ║"
echo "║                  (Jest-like Output)                            ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Run the tests with verbose output
cd "$(dirname "$0")"
go test -v ./internal/routes/... -count=1 2>&1

echo ""
