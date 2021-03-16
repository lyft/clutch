SHELL:=/usr/bin/env bash
.DEFAULT_GOAL:=all

MAKEFLAGS += --no-print-directory

PRIMARY_MODULE_DIR := $(shell cd backend && go list -f "{{ "{{" }} .Dir {{ "}}" }}" -m "github.com/lyft/clutch/backend")
TOOLS_MODULE_DIR := $(shell cd backend && go list -f "{{ "{{" }} .Dir {{ "}}" }}" -m "github.com/lyft/clutch/tools")
MY_ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: all # Generate API, Frontend, and backend assets.
all: api frontend backend-tidy backend-with-assets

.PHONY: api
api: yarn-ensure
	$(SHELL) $(TOOLS_MODULE_DIR)/compile-protos.sh $(MY_ROOT_DIR)

.PHONY: backend
backend:
	cd backend && go build -o ../build/clutch main.go

.PHONY: backend-dev
backend-dev:
	REPO_ROOT=$(MY_ROOT_DIR) $(SHELL) $(TOOLS_MODULE_DIR)/air.sh

.PHONY: backend-test
backend-test:
	cd backend && go test -race -covermode=atomic ./...

backend-tidy:
	cd backend && go mod tidy

.PHONY: backend-with-assets
backend-with-assets: frontend
	cd backend && go run $(PRIMARY_MODULE_DIR)/cmd/assets/generate.go ../frontend/build && go build -tags withAssets -o ../build/clutch -ldflags="-X main.version=$(VERSION)"

.PHONY: frontend
frontend: yarn-install
	yarn --cwd frontend build

.PHONY: yarn-install
yarn-install: yarn-ensure
	yarn --cwd frontend install

.PHONY: yarn-ensure
yarn-ensure:
	@$(SHELL) $(TOOLS_MODULE_DIR)/install-yarn.sh

.PHONY: frontend-dev # Start the frontend in development mode.
frontend-dev: yarn-ensure
	yarn --cwd frontend install --frozen-lockfile && yarn --cwd frontend start
