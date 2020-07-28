SHELL:=/usr/bin/env bash
.DEFAULT_GOAL:=all

MAKEFLAGS += --no-print-directory

PRIMARY_MODULE_DIR := $(shell cd backend && go list -f "{{ "{{" }} .Dir {{ "}}" }}" -m "github.com/lyft/clutch/backend")
TOOLS_MODULE_DIR := $(shell cd backend && go list -f "{{ "{{" }} .Dir {{ "}}" }}" -m "github.com/lyft/clutch/tools")
MY_ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: all # Generate API, Frontend, and backend assets.
all: api frontend backend-with-assets

.PHONY: api
api: yarn-ensure
	$(SHELL) $(TOOLS_MODULE_DIR)/compile-protos.sh $(MY_ROOT_DIR)

.PHONY: backend
backend:
	cd backend && go build -o ../build/clutch main.go

.PHONY: backend-dev
backend-dev:
	cd backend && go run main.go -c clutch-config.yaml

.PHONY: backend-test
backend-test:
	cd backend && go test -race -covermode=atomic ./...

.PHONY: backend-with-assets
backend-with-assets: frontend
	cd backend && go run $(PRIMARY_MODULE_DIR)/cmd/assets/generate.go ../frontend/build && go build -tags withAssets -o ../build/clutch -ldflags="-X main.version=$(VERSION)"

.PHONY: vendor
vendor:
	mkdir -p build/tools
	cd backend && go get github.com/lyft/clutch/backend@latest

frontend: yarn-install
	yarn --cwd frontend build

.PHONY: yarn-install
yarn-install: yarn-ensure
	yarn --cwd frontend install

.PHONY: yarn-ensure
yarn-ensure:
	@$(SHELL) $(TOOLS_MODULE_DIR)/install-yarn.sh
