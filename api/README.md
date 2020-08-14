# API Guidelines

- Always use `proto3`. [Google Language Guide](https://developers.google.com/protocol-buffers/docs/proto3)
- Always use a distinct message for RPC requests and responses.

https://medium.com/@akhaku/protobuf-definition-best-practices-87f281576f31

## Required Annotations

## grpc-gateway

Clutch uses [grpc-ecosystem/grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) to serve endpoints via HTTP+JSON for the front-end in addition to the gRPC normally provided by the proto definitions.

grpc-gateway uses the [`google.api.http` annotation](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto#L44-L312) to map endpoints.

Rules for `google.api.http` annotations:
- Always use `POST`.
- Don't forget to include `body : "*"` or data will not be transmitted.
- The API endpoint should always be camel case.
- The API endpoint should be literally named after the package and RPC method, e.g. `package clutch.aws.ec2.v1`

##  Enums

All `enum` values should have a zero value of `UNSPECIFIED`. This signifies that the value was never filled in.
If a real value was used for the zero value, then that value may be erroneously presented in an object despite not being
filled in.

`UNKNOWN` is also a useful value to be used when the value needed is not present in the enum for some reason, for example when
translating between two systems where a string will become an enumerated value.

```proto
enum Foo {
  UNSPECIFIED = 0;  // Value was not filled in.
  UNKNOWN = 1;      // Value could not be represented by existing enum values.
  MY_VALUE = 2;
  MY_OTHER_VALUE = 3;
}
```
