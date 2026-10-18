#!/bin/bash
set -e
echo "Formatting Go packages..."
go fmt ./cmd/... ./client/golang/... ./server/... ./pkgs/... ./test/...
