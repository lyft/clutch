FROM golang:1.16.0

ENTRYPOINT cd /code/module/chaos/serverexperimentation/xds && go test -tags integration_only -c -o /testrunner
