FROM alpine

ADD testrunner /testrunner

ENTRYPOINT /testrunner
