FROM golang:1.16.0

COPY . /code

RUN cd /code/module/chaos/serverexperimentation && find . -type f -name "*_test.go" ! -name "*_integration_test.go" -exec rm {} \;
RUN cd /code/module/chaos/serverexperimentation && go test -tags integration_only -c -o /testrunner
