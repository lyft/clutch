# API Guidelines

<!-- TO UPDATE ToC: run `npx doctoc README.md`>
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [General](#general)
- [Enums](#enums)
- [Rules](#rules)
  - [All HTTP annotations should use the POST method](#all-http-annotations-should-use-the-post-method)
  - [Version should always be the trailing component of package](#version-should-always-be-the-trailing-component-of-package)
  - [Package name and folder structure should match](#package-name-and-folder-structure-should-match)
  - [Service name should match last non-version component of package name](#service-name-should-match-last-non-version-component-of-package-name)
  - [HTTP annotation path should match package name, annotation should be same as RPC method](#http-annotation-path-should-match-package-name-annotation-should-be-same-as-rpc-method)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## General

- Always use `proto3`. [Google Language Guide](https://developers.google.com/protocol-buffers/docs/proto3)
- Always use a distinct message for RPC requests and responses.

For more general reading on proto best practices see https://medium.com/@akhaku/protobuf-definition-best-practices-87f281576f31

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

## Rules

Below are conventions followed in the Clutch APIs that are unfortunately not automatically linted yet.

The rules are meant to take as much guesswork out of naming things as possible by making the filename, package name, and services match. HTTP annotations are literal RPC, instead of trying to map RESTful concepts back onto the API.

### All HTTP annotations should use the POST method

**Examples of incorrect definition for this rule:**
```
rpc GetThing(GetThingRequest) returns (GetThingResponse) {
  option (google.api.http) = {
    get : "/v1/foo/thing"
  };
}
```
- :x: HTTP annotation uses GET instead of POST.

**Examples of correct definition for this rule:**
```
rpc GetThing(GetThingRequest) returns (GetThingResponse) {
  option (google.api.http) = {
    post : "/v1/foo/getThing"
    body : "*"
  };
}
```

### Version should always be the trailing component of package

**Examples of incorrect definition for this rule:**
```proto
package clutch.aws.v1.ec2;
```

- :x: Last component of package is `ec2`, not `v1`.

```proto
package clutch.v1.healthcheck;
```

- :x: Last component of package is `healthcheck`, not `v1`.

**Examples of correct definition for this rule:**
```proto
package clutch.healthcheck.v1;
```

```proto
package clutch.aws.ec2.v1;
```

### Package name and folder structure should match

**Examples of incorrect definition for this rule:**

- Filename `clutch/api/aws/v1/ec2.proto`
```proto
package clutch.aws.ec2.v1;
```

- :x: `ec2` is missing from folder structure.

**Examples of correct definition for this rule:**

- Filename `clutch/api/aws/ec2/v1/ec2.proto`
```proto
package clutch.aws.ec2.v1;
```

### Service name should match last non-version component of package name

**Examples of incorrect definition for this rule:**

```proto
package clutch.foo.v1;

service BarAPI {
  ...
}
```

- :x: Last package component is `foo`, but service name is `bar`.

**Examples of correct definition for this rule:**

- Filename `clutch/api/foo/v1/foo.proto`
```proto
package clutch.foo.v1;

service FooAPI {
  ...
}
```

- Filename `clutch/api/foo/bar/v1/bar.proto`
```proto
package clutch.foo.bar.v1;

service BarAPI {
  ...
}
```

### HTTP annotation path should match package name, annotation should be same as RPC method

**Examples of incorrect definition for this rule:**

```proto
package clutch.foo.v1;

service FooAPI {
  rpc CreateBar(CreateBarRequest) returns (CreateBarResponse) {
    option (google.api.http) = {
      post : "/v1/bar/create"
      body : "*"
    };
  }
}
```

- :x: Leading path is `/v1/experiments`, not `/v1/foo` in accordance with package name.
- :x: RPC method is `CreateExperiment` but HTTP mapping is `/experiments/create`.

**Examples of correct definition for this rule:**

```proto
package clutch.foo.v1;

service FooAPI {
  rpc CreateBar(CreateBarRequest) returns (CreateBarResponse) {
    option (google.api.http) = {
      post : "/v1/foo/createBar"
      body : "*"
    };
  }
}
```

