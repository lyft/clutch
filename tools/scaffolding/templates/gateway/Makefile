SHELL:=/usr/bin/env bash
.DEFAULT_GOAL:=all

MAKEFLAGS += --no-print-directory

PRIMARY_MODULE_DIR := $(shell cd backend && go list -f "{{ "{{" }} .Dir {{ "}}" }}" -m "github.com/lyft/clutch/backend")
MY_ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: all # Generate API, Frontend, and backend assets.
all: api frontend backend-with-assets

.PHONY: api
api: yarn-ensure
	$(PRIMARY_MODULE_DIR)/../tools/compile-protos.sh $(MY_ROOT_DIR)

.PHONY: backend
backend:
	cd backend && go build -o ../build/clutch main.go

.PHONY: backend-dev
backend-dev:
	cd backend && go run main.go -c ../clutch-config.yaml

.PHONY: backend-test
backend-test:
	cd backend && go test -race -covermode=atomic ./...

.PHONY: backend-with-assets
backend-with-assets: frontend
	cd backend && go run $(PRIMARY_MODULE_DIR)/cmd/assets/generate.go ../frontend/build && go build -tags withAssets -o ../build/clutch -ldflags="-X main.version=$(VERSION)"

.PHONY: frontend
frontend: yarn-ensure
	$(PRIMARY_MODULE_DIR)/../tools/frontend-gateway-linker.sh $(MY_ROOT_DIR) build

.PHONY: yarn-ensure
yarn-ensure:
	@$(PRIMARY_MODULE_DIR)/../tools/install-yarn.sh
