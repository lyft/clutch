FROM golang:1.16.0

COPY . /code

RUN cd /code/module/chaos/serverexperimentation/xds/integration && go test -tags integration_only -c -o /testrunner
