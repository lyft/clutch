target_dir="$(pwd)"

export GOOS=linux
export GOARCH=amd64

pushd ../../../..
	pushd internal/test/integration/xds/cmd/envoyconfiggen
	  go build -o $target_dir/envoyconfiggen main.go
	popd
	pushd module/chaos/serverexperimentation/xds
		go test -tags integration_only -c -o $target_dir/testrunner
	popd
popd

docker-compose up --build --abort-on-container-exit
