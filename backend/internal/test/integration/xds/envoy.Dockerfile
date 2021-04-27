FROM golang:1.16.0 AS build
COPY envoyconfiggen /envoyconfiggen
RUN /envoyconfiggen > /config.json

FROM envoyproxy/envoy:v1.17.1
COPY --from=build /config.json /config.json

ENTRYPOINT /usr/local/bin/envoy -c /config.json
