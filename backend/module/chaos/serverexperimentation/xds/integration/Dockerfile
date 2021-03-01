FROM golang:1.15.6

COPY . /build

RUN cd /build/module/chaos/serverexperimentation/xds/integration && go test -tags integration_only -c -o /test
