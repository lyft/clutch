---
title: Mock Gateway
{{ .EditURL }}
---

The mock gateway and mock components, located in `backend/mock` allow testing and development of features without relying on dependencies such as third-party providers. This is especially useful for fine-tuning the UI or demoing features that include destructive capabilities.

## Services

In order for a component to be used in the mock server, there must be a corresponding mock implementation.

Mock services can be found in [`backend/mock/service`](https://github.com/lyft/clutch/blob/main/backend/mock/service).

## Gateway

The mock gateway consumes a regular Clutch backend config with mock components registered instead of the real components. See [`backend/mock/gateway.go`](https://github.com/lyft/clutch/blob/main/backend/mock/gateway.go).

This is most commonly done in the form of mocked services since all external dependencies are generally hidden behind services.

It's also possible to run with a mix of real and mocked components by including both in the list of components passed to the gateway's `Run` function.

## Requirements

See [Local Build Requirements](/docs/getting-started/local-build#requirements).

## Running the Mock Gateway

`make dev-mock` will compile and start the mock server on `localhost:8080` and a hot-reloading frontend on `localhost:3000`.
