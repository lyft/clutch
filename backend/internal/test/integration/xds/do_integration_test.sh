#!/usr/bin/env bash
set -euo pipefail

# We expect this to be run from the directory containing the script.
target_dir="$(pwd)"

# Build the Go executable on the host system targeting linux. This is much faster
# than naively building within the container setup since we get to share the host Go
# caches.
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# Place the output binaries in $target_dir for use by the Dockerfiles.
pushd ../../../..
	pushd internal/test/integration/xds/cmd/envoyconfiggen
	  go build -o $target_dir/envoyconfiggen main.go
	popd
	pushd module/chaos/serverexperimentation/xds
		go test -tags integration_only -c -o $target_dir/testrunner
	popd
popd

docker-compose up --build --abort-on-container-exit
