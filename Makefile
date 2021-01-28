SHELL:=/usr/bin/env bash
.DEFAULT_GOAL:=all

MAKEFLAGS += --no-print-directory

DOCS_DEPLOY_USE_SSH ?= true
DOCS_DEPLOY_GIT_USER ?= git

VERSION := 0.0.0

YARN:=./build/bin/yarn.sh
PROJECT_ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: help # Print this help message.
 help:
	@grep -E '^\.PHONY: [a-zA-Z_-]+ .*?# .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = "(: |#)"}; {printf "%-30s %s\n", $$2, $$3}'

.PHONY: all # Generate API, Frontend, and backend assets.
all: api frontend backend-with-assets

.PHONY: api # Generate API assets.
api: yarn-ensure
	tools/compile-protos.sh -c "$(PROJECT_ROOT_DIR)/api"

.PHONY: api-lint # Lint the generated API assets.
api-lint:
	tools/compile-protos.sh -c "$(PROJECT_ROOT_DIR)/api" -l

.PHONY: api-lint-fix # Lint and fix the generated API assets.
api-lint-fix:
	tools/compile-protos.sh -c "$(PROJECT_ROOT_DIR)/api" -lf

.PHONY: api-verify # Verify API proto changes include generate frontend and backend assets.
api-verify:
	find backend/api -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} \;
	find frontend/api/src -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} \;
	$(MAKE) api
	tools/ensure-no-diff.sh backend/api frontend/api/src

.PHONY: backend # Build the standalone backend.
backend:
	cd backend && go build -o ../build/clutch -ldflags="-X main.version=$(VERSION)"

.PHONY: backend-with-assets # Build the backend with frontend assets.
backend-with-assets:
	cd backend && go run cmd/assets/generate.go ../frontend/packages/app/build && go build -tags withAssets -o ../build/clutch -ldflags="-X main.version=$(VERSION)"

.PHONY: backend-dev # Start the backend in development mode.
backend-dev:
	tools/air.sh

.PHONY: backend-dev-mock # Start the backend in development mode with mock responses.
backend-dev-mock:
	cd backend && go run mock/gateway.go

.PHONY: backend-lint # Lint the backend code.
backend-lint:
	tools/golangci-lint.sh run

.PHONY: backend-lint-fix # Lint and fix the backend code.
backend-lint-fix:
	tools/golangci-lint.sh run --fix
	cd backend && go mod tidy

.PHONY: backend-test # Run unit tests for the backend code.
backend-test:
	cd backend && go test -race -covermode=atomic ./...

.PHONY: backend-verify # Verify go modules' requirements files are clean.
backend-verify:
	cd backend && go mod tidy
	tools/ensure-no-diff.sh backend

.PHONY: backend-config-validation
backend-config-validation:
	cd backend && go run main.go -validate -c clutch-config.yaml

.PHONY: yarn-install # Install frontend dependencies.
yarn-install: yarn-ensure
	$(YARN) --cwd frontend install --frozen-lockfile 

.PHONY: backend-integration-test
backend-integration-test:
	cd backend/module/chaos/serverexperimentation/rtds/integration && docker-compose up --build --abort-on-container-exit

.PHONY: frontend # Build production frontend assets.
frontend: yarn-install
	$(YARN) --cwd frontend build

.PHONY: frontend-dev-build # Build development frontend assets.
frontend-dev-build: yarn-install
	$(YARN) --cwd frontend build:dev

.PHONY: frontend-dev # Start the frontend in development mode.
frontend-dev: yarn-install
	$(YARN) --cwd frontend start

.PHONY: frontend-lint # Lint the frontend code.
frontend-lint: yarn-ensure
	$(YARN) --cwd frontend lint

.PHONY: frontend-lint-fix # Lint and fix the frontend code.
frontend-lint-fix: yarn-ensure
	$(YARN) --cwd frontend lint:fix

.PHONY: frontend-test # Run unit tests for the frontend code.
frontend-test: yarn-ensure
	$(YARN) --cwd frontend test

.PHONY: frontend-e2e # Run end-to-end tests for the frontend code.
frontend-e2e: yarn-ensure
	./tools/frontend-e2e.sh

.PHONY: frontend-verify # Verify frontend packages are sorted.
frontend-verify: yarn-ensure
	$(YARN) --cwd frontend lint:packages

.PHONY: docs # Build all doc assets.
docs: docs-generate yarn-ensure
	$(YARN) --cwd docs/_website install --frozen-lockfile && $(YARN) --cwd docs/_website build

.PHONY: docs-dev # Start the docs server in development mode.
docs-dev: docs-generate yarn-ensure
	$(YARN) --cwd docs/_website install --frozen-lockfile && BROWSER=none $(YARN) --cwd docs/_website start

.PHONY: docs-generate # Generate the documentation content.
docs-generate:
	cd docs/_website/generator && go run .

.PHONY: dev # Run the Clutch application in development mode.
dev:
	$(MAKE) -j2 backend-dev frontend-dev

.PHONY: dev-mock # Run the Clutch application in development mode with mock responses.
dev-mock:
	$(MAKE) -j2 backend-dev-mock frontend-dev

.PHONY: lint # Lint all of the code.
lint: api-lint backend-lint frontend-lint

.PHONY: lint-fix # Lint and fix all of the code.
lint-fix: api-lint-fix backend-lint-fix frontend-lint-fix

.PHONY: scaffold-gateway # Generate a new gateway.
scaffold-gateway:
	cd tools/scaffolding && go run scaffolder.go -m gateway -p $(shell git rev-parse --short HEAD)

.PHONY: scaffold-workflow # Generate a new Workflow package.
scaffold-workflow:
	cd tools/scaffolding && go run scaffolder.go -m frontend-plugin

.PHONY: storybook # Start storybook locally.
storybook: yarn-install
	$(YARN) --cwd frontend storybook

.PHONY: storybook-build # Build storybook assets for deploy.
storybook-build: yarn-install
	$(YARN) --cwd frontend storybook:build

.PHONY: test # Unit test all of the code.
test: backend-test frontend-test

.PHONY: verify # Verify all of the code.
verify: api-verify backend-verify frontend-verify

.PHONY: yarn-ensure # Install the pinned version of yarn.
yarn-ensure:
	@./tools/install-yarn.sh

.PHONY: dev-k8s-up # Start a local k8s cluster
dev-k8s-up:
	@tools/kind.sh create cluster --kubeconfig $(PROJECT_ROOT_DIR)/build/kubeconfig-clutch --name clutch-local || true
	@tools/kind.sh seed

	@echo
	@echo "Export these environment variables before starting development:"
	@echo '    export KUBECONFIG=$(PROJECT_ROOT_DIR)/build/kubeconfig-clutch'

.PHONY: k8s-stop
dev-k8s-down:
	tools/kind.sh delete cluster --name clutch-local

