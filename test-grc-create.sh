#!/bin/bash
set -e

# Create Project
cd ..
rm -rf testing-build
mkdir testing-build
cd testing-build
go run ./../go-router/grc/main.go -cv $( git log -1 origin/$(git rev-parse --abbrev-ref HEAD) | grep commit | sed 's/commit//g' > git-log.txt) create github.com/jetnoli/testing
cd testing

# Asserts
go run main.go