#!/bin/bash
set -e

# Create Project
goRouterVersion=$(git log -1 origin/$(git rev-parse --abbrev-ref HEAD) | grep commit | sed 's/commit//g' | sed 's/^[[:space:]]*//; s/[[:space:]]*$//' | tr '\n' ' ')

cd ../../testing-build
rm -rf page_map
mkdir page_map
cd page_map
go mod init app-builder
go get github.com/jetnoli/go-router@${goRouterVersion}
cat ../../go-router/grc/_go.txt > main.go
go mod tidy
go run main.go ../../go-router/grc/static/