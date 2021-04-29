FROM alpine

ADD build/testrunner /testrunner

ENTRYPOINT /testrunner
