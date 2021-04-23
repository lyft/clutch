FROM golang:1.16.0 AS build
COPY . /code
RUN cd /code/internal/test/integration/serverexperimentation/cmd/envoyconfiggen && go build -o /tmp/generate_config
RUN /tmp/generate_config > /config.json

FROM envoyproxy/envoy:v1.17.1
COPY --from=build /config.json /config.json

ENTRYPOINT /usr/local/bin/envoy -c /config.json
