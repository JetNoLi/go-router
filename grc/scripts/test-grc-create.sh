#!/bin/bash
set -e

# Create Project
goRouterVersion=$(git log -1 origin/$(git rev-parse --abbrev-ref HEAD) | grep commit | sed 's/commit//g' | sed 's/^[[:space:]]*//; s/[[:space:]]*$//' | tr '\n' ' ')
cd ../
rm -rf testing-build
mkdir testing-build
cd testing-build
echo "version $goRouterVersion"
echo "$(ls ../go-router/grc)"
cp ../go-router/grc/go.mod ./
cp ../go-router/grc/go.sum ./
go run ../go-router/grc/main.go -cv "$goRouterVersion" create github.com/jetnoli/testing
cd testing

# Asserts
go run main.go