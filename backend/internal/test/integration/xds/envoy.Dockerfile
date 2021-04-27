
FROM envoyproxy/envoy-alpine:v1.17.1 as envoy

FROM golang:1.16.0
COPY --from=envoy /usr/local/bin/envoy /usr/local/bin/envoy

ENV GO111MODULE=on

ENTRYPOINT bash -c '/usr/local/bin/envoy --config-yaml "$(cd /code/internal/test/integration/xds/cmd/envoyconfiggen && go run main.go)"'
