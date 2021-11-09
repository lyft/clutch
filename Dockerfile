# Frontend build.
FROM node:16.13.0-buster as nodebuild
COPY ./frontend ./frontend
COPY ./tools/install-yarn.sh ./tools/install-yarn.sh
COPY ./tools/preflight-checks.sh ./tools/preflight-checks.sh
COPY Makefile .

RUN make frontend

# Backend build.
FROM golang:1.17.3-buster as gobuild
WORKDIR /go/src/github.com/lyft/clutch
COPY ./backend ./backend
COPY ./tools/preflight-checks.sh ./tools/preflight-checks.sh
COPY Makefile .

COPY --from=nodebuild ./frontend/packages/app/build ./frontend/packages/app/build/

RUN make backend-with-assets

# Copy binary to final image.
FROM gcr.io/distroless/base-debian10
COPY --from=gobuild /go/src/github.com/lyft/clutch/build/clutch /
COPY backend/clutch-config.yaml /
CMD ["/clutch"]
