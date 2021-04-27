
FROM envoyproxy/envoy:v1.17.1 as envoy

FROM golang:1.16.0
COPY --from=envoy /usr/local/bin/envoy /usr/local/bin/envoy

ENTRYPOINT bash -c '/usr/local/bin/envoy -c $(go run /code/internal/test/integration/xds/cmd/envoyconfiggen)'
