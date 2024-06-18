FROM envoyproxy/envoy:v1.30.0

COPY build/envoyconfiggen /envoyconfiggen
RUN /envoyconfiggen > /config.json

ENTRYPOINT /usr/local/bin/envoy -c /config.json
