SHELL:=/usr/bin/env bash
.DEFAULT_GOAL:=all

MAKEFLAGS += --no-print-directory

DOCS_DEPLOY_USE_SSH ?= true
DOCS_DEPLOY_GIT_USER ?= git

VERSION := 0.0.0

.PHONY: help # Print this help message.
 help:
	@grep -E '^\.PHONY: [a-zA-Z_-]+ .*?# .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = "(: |#)"}; {printf "%-30s %s\n", $$2, $$3}'

.PHONY: all # Generate API, Frontend, and backend assets.
all: yarn-ensure api frontend backend-with-assets

.PHONY: api # Generate API assets.
api:
	tools/compile-protos.sh

.PHONY: api-lint # Lint the generated API assets.
api-lint:
	tools/compile-protos.sh -l

.PHONY: api-lint-fix # Lint and fix the generated API assets.
api-lint-fix:
	tools/compile-protos.sh -lf

.PHONY: api-verify # Verify API proto changes include generate frontend and backend assets.
api-verify:
	find backend/api -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} \;
	find frontend/api -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} \;
	$(MAKE) api
	tools/ensure-no-diff.sh backend/api frontend/api

.PHONY: backend # Build the standalone backend.
backend:
	cd backend && go build -o ../build/clutch -ldflags="-X main.version=$(VERSION)"

.PHONY: backend-with-assets # Build the backend with frontend assets.
backend-with-assets:
	cd backend && go run cmd/assets/generate.go ../frontend/packages/app/build && go build -tags withAssets -o ../build/clutch -ldflags="-X main.version=$(VERSION)"

.PHONY: backend-dev # Start the backend in development mode.
backend-dev:
	cd backend && go run .

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

.PHONY: frontend # Build production frontend assets.
frontend: yarn-ensure
	cd frontend && yarn install --frozen-lockfile && yarn build

.PHONY: frontend-dev-build # Build development frontend assets.
frontend-dev-build: yarn-ensure
	cd frontend && yarn install --frozen-lockfile && yarn build:dev

.PHONY: frontend-dev # Start the frontend in development mode.
frontend-dev: yarn-ensure
	cd frontend && yarn install --frozen-lockfile && yarn start

.PHONY: frontend-lint # Lint the frontend code.
frontend-lint: yarn-ensure
	cd frontend && yarn lint

.PHONY: frontend-lint-fix # Lint and fix the frontend code.
frontend-lint-fix: yarn-ensure
	cd frontend && yarn lint:fix

.PHONY: frontend-test # Run unit tests for the frontend code.
frontend-test: yarn-ensure
	cd frontend && yarn test

.PHONY: frontend-e2e # Run end-to-end tests for the frontend code.
frontend-e2e: yarn-ensure
	./tools/frontend-e2e.sh

.PHONY: frontend-verify # Verify frontend packages are sorted.
frontend-verify: yarn-ensure
	cd frontend && yarn lint:packages

.PHONY: docs # Build all doc assets.
docs: docs-generate
	cd docs/_website && yarn install --frozen-lockfile && yarn build

.PHONY: docs-deploy # Deploy the documentation.
docs-deploy: docs
	cd docs/_website && GIT_USER=$(DOCS_DEPLOY_GIT_USER) USE_SSH=$(DOCS_DEPLOY_USE_SSH) yarn deploy

.PHONY: docs-dev # Start the docs server in development mode.
docs-dev: docs-generate
	cd docs/_website && yarn install --frozen-lockfile && BROWSER=none yarn start

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
	cd tools/scaffolding && go run scaffolder.go -m gateway

.PHONY: scaffold-workflow # Generate a new Workflow package.
scaffold-workflow:
	cd tools/scaffolding && go run scaffolder.go -m frontend-plugin

.PHONY: test # Unit test all of the code.
test: backend-test frontend-test

.PHONY: verify # Verify all of the code.
verify: api-verify backend-verify frontend-verify

.PHONY: yarn-ensure # Install the pinned version of yarn.
yarn-ensure:
	@./tools/install-yarn.sh
