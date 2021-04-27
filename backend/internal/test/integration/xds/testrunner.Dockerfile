FROM golang:1.16.0

COPY . /code

#RUN envoy -c $(go run whatever/main.go)
RUN go test ./experimentation/xds -tags integration_only
#RUN cd /code/module/chaos/serverexperimentation && go test -tags integration_only -c -o /testrunner
