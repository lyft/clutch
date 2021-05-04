#!/usr/bin/env bash
set -euo pipefail

# We expect this to be run from the directory containing the script.
mkdir -p build/
target_dir="$(pwd)/build"

# Build the Go executable on the host system targeting linux. This is much faster
# than naively building within the container setup since we get to share the host Go
# caches.
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# Run integration tests from a provided directory.
#
# @arg $1 path to an Envoy config generator file called `main.go` that should be used
# to generate a config that's passed to an Envoy instance that is spin up as part of
# the process of running integration tests.
# @arg $2 path to a directory with integration tests to run.
run_tests () {
	# Place the output binaries in $target_dir for use by the Dockerfiles.
	pushd ../../../..
		pushd $1
	  		go build -o $target_dir/envoyconfiggen main.go
		popd
		pushd $2
			go test -tags integration_only -c -o $target_dir/testrunner
		popd
  	popd

	docker-compose up --build --abort-on-container-exit
}

run_tests "internal/test/integration/chaos/experimentation/cmd/envoyconfiggen" "module/chaos/experimentation/xds"
run_tests "internal/test/integration/chaos/serverexperimentation/cmd/envoyconfiggen" "module/chaos/serverexperimentation"
