---
title: API Definitions
{{ .EditURL }}
---

Clutch uses `proto3` to define API endpoints and objects in addition to backend configuration schemas. Protobuf is an interface definition language (IDL) from Google. See the [Protocol Buffer's Language Guide](https://developers.google.com/protocol-buffers/docs/proto) for more information.

Protobuf was chosen because it has a rich tooling ecosystem, and can generate clients or server implementation stubs for nearly any target language. Clutch uses protobuf tooling for:
- JSON <-> gRPC transcoding from a single API definition
- Server API stub implementations and API objects
- Frontend API objects
- Input validation code from annotations
- Auditing and authorization configuration for each API endpoint from annotations

Protobuf also allows you to annotate APIs and objects with additional metadata, which can then be used for additional code generation, or read at runtime for configuration purposes. For example, Clutch uses [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) to provide a single server that can respond to requests with either gRPC/JSON and [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate) for input validation.

## Structure

Clutch protobuf definitions are in the [`api/`](https://github.com/lyft/clutch/blob/main/api) directory.

```
├── api        # clutch-specific annotation definitions
├── config     # backend config schemas for the gateway, services, etc.
└── buf.json   # api linter configuration   
```

### Modules

Other directories in the top-level `api/` folder correspond to API definitions for modules. For example, [`api/healthcheck/`](https://github.com/lyft/clutch/blob/main/api/healthcheck/v1/healthcheck.proto).

### Versioning and Backward Compatibility

Directories should always be named with the version number `v1` as the last component.

While Clutch is still in its early stages, we make no explicit guarantees about backwards compatibility. However, we will avoid making breaking changes if at all possible, and publicize them in a changelog.

Protobuf has tooling to prevent breaking changes that will be adopted at a later date.

## Generating Code

From the root of the project directory, run:

```bash
make api
```

Generated code from the API definitions lives in `backend/api` and `frontend/api/src`. Files in these directories should never be edited directly.

:::info Note
The [script that compiles protos](https://github.com/lyft/clutch/blob/main/tools/compile-protos.sh) automatically downloads `protoc`, `pbjs`, and other dependencies to the local build environment.
:::

## JSON Endpoint Availability

By default, the protobuf toolchain will only generate gRPC stubs for API definitions. The `google.api.http` annotation is used to also make the endpoint available via the JSON server. 

For an example, see [`api/healthcheck/v1/healthcheck.proto`](https://github.com/lyft/clutch/blob/main/api/healthcheck/v1/healthcheck.proto).

For more information see, [grpc-ecosystem/grpc-gateway#usage](https://github.com/grpc-ecosystem/grpc-gateway#usage).

:::info Method Information
Clutch registers all API endpoints (other than healthcheck) with the `POST` method to more easily map back to the native RPC semantics of protobuf.
:::

Always register endpoints

## Clutch-specific Annotations

TODO: document Clutch-specific annotations

## API Guidelines

Clutch uses [`buf`](https://buf.build/docs/introduction) for linting and `clang-format` for formatting protobufs. Run the following command to automatically fix any fixable issues:

```bash
make api-lint-fix
```

Rules and recommendations:
- All `enum` values should have a zero value of `UNSPECIFIED`. `UNKNOWN` should also be included if the value is being translated from another system and may change.
- Always use distinct message types for RPC requests and responses.
- Always use `POST` for the JSON method.
- Always name endpoints and request/response types literally after their proto RPC definition, e.g. `FetchConfig(FetchConfigRequest) returns (FetchConfigResponse)` which maps to `post : "/v1/config/fetchConfig"`.
- Using request/response protobuf objects all the way down to the service is encouraged as there is no benefit to an intermediate object if the information required is identical.

## The `Any` Type
Clutch makes heavy use of the `Any` type in protobuf for managing component configuration. This makes it easy to dynamically handle whatever configuration types were available at build  time, whether from the core project or an extension. The tradeoff is that any type mismatch when attempting to deserialize will throw an error at runtime. We currently only attempt deserialization when instantiating components in the gateway, so the gateway will hard fail on start in the event a type is not found.

:::info
`typed_config` is the nominal name for fields with the `Any` type in Clutch configuration.
:::

## Database Objects

While it is tempting to use protobuf all the way from the edge to the database, it is recommended to have separate objects and schemas to serialize data to a persistent store. Updating the API and migrating the database at the same time is a tricky proposition.
