FROM golang:1.16.0 AS build
COPY . /code
RUN cd /code/internal/test/integration/xds/cmd/envoyconfiggen && go build -o /generate_config
