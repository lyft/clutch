# Frontend build.
FROM node:15-buster as nodebuild
COPY ./frontend ./frontend
COPY ./tools/install-yarn.sh ./tools/install-yarn.sh
COPY Makefile .

RUN make frontend

# Backend build.
FROM golang:1.16-buster as gobuild
WORKDIR /go/src/github.com/lyft/clutch
COPY ./backend ./backend
COPY Makefile .

COPY --from=nodebuild ./frontend/packages/app/build ./frontend/packages/app/build/

RUN make backend-with-assets

# Copy binary to final image.
FROM gcr.io/distroless/base-debian10
COPY --from=gobuild /go/src/github.com/lyft/clutch/build/clutch /
COPY backend/clutch-config.yaml /
CMD ["/clutch"]
