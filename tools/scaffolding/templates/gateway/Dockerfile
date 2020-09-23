FROM node:14-buster as nodebuild
COPY ./frontend ./frontend
COPY ./tools/install-yarn.sh ./tools/install-yarn.sh
COPY Makefile .
RUN tools/install-yarn.sh && yarn --cwd frontend install

FROM golang:1.15-buster as gobuild
WORKDIR /go/src/{{ .RepoProvider}}/{{ .RepoOwner }}/{{ .RepoName }}
COPY ./backend ./backend
COPY Makefile .
COPY --from=nodebuild ./frontend/build ./frontend/build/
RUN cd backend \
    && go mod download \
    && go run $(go list -f "{{ .Dir }}" -m "github.com/lyft/clutch/backend")/cmd/assets/generate.go ../frontend/build \
    && go build -tags withAssets -o ../build/clutch -ldflags="-X main.version=${VERSION}"

FROM gcr.io/distroless/base-debian10
COPY --from=gobuild /go/src/{{ .RepoProvider}}/{{ .RepoOwner }}/{{ .RepoName }}/build/clutch /
COPY backend/clutch-config.yaml /
CMD ["/clutch"]
