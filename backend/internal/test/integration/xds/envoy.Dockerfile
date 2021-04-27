FROM envoyproxy/envoy-alpine:v1.17.1

COPY envoyconfiggen /envoyconfiggen
RUN /envoyconfiggen > /config.json

ENTRYPOINT /usr/local/bin/envoy -c /config.json
