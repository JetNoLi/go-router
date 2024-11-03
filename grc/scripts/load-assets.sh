#!/bin/bash
set -e

# Create Project
goRouterVersion=$(git log -1 origin/$(git rev-parse --abbrev-ref HEAD) | grep commit | sed 's/commit//g' | sed 's/^[[:space:]]*//; s/[[:space:]]*$//' | tr '\n' ' ')

# cd ../../testing-build
rm -rf page_map
mkdir page_map
cd page_map
go mod init app-builder
go get github.com/jetnoli/go-router@${goRouterVersion}
cat ../_go.txt > main.go
go mod tidy
go run main.go ../static

#!/bin/bash
set -e

#TODO: Turn in to grc command

# # Create Project
# goRouterVersion=d0fc733439be19e4cd044869da84baf7eaeb1ce6

# # cd ../../testing-build
# rm -rf page_map
# mkdir page_map
# cd page_map
# go mod init app-builder
# go get github.com/jetnoli/go-router@${goRouterVersion}
# cat ../_go.txt > main.go
# go mod tidy
# cd ..
# go run page_map/main.go ./
# rm -rf page_map