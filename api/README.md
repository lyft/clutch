# API Guidelines

- Always use `proto3`. [Google Language Guide](https://developers.google.com/protocol-buffers/docs/proto3)
- Always use a distinct message for RPC requests and responses.

https://medium.com/@akhaku/protobuf-definition-best-practices-87f281576f31

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