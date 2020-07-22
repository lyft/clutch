---
title: Configuration Guide
{{ .EditURL }}
---

## Backend

### Command-line Arguments

The Clutch binary will consume `clutch-config.yaml` in the current working directory by default.

To use a different path for the configuration,use the `-c` option, e.g.

```bash
./clutch -c /etc/clutch/clutch-config.yaml
```

### Config Features

Clutch supports the expansion of environment variables after reading the YAML when the gateway starts up.

```yaml
password: ${MY_SECRET_PASSWORD}
```

### YAML Specification

All backend configuration in Clutch is specified in protobuf definitions. For information on how YAML and JSON map to protobuf see [Language Guide (proto3): JSON Mapping](https://developers.google.com/protocol-buffers/docs/proto3#json).

*Note: it is recommended to write YAML long-form, the format of the examples below are shortened for documentation purposes.*

#### Gateway
See [`api/config/gateway/v1/gateway.proto`](https://github.com/lyft/clutch/blob/main/api/config/gateway/v1/gateway.proto) for the full configuration specification. For an example of a filled-in config, see the [sample `clutch-config.yaml`](https://github.com/lyft/clutch/blob/main/backend/clutch-config.yaml).

##### Top-level Configuration
```yaml title="clutch-config.yaml"
{{ simpleProtoYAML "clutch.config.gateway.v1.Config" }}
```

##### `GatewayOptions`
```yaml
{{ simpleProtoYAML "clutch.config.gateway.v1.GatewayOptions" }}
```

##### `Module`, `Resolver`, `Service`
Modules, resolvers, and service are all specified using the same format. The [name of the component](/docs/components#backend) is specified, and if necessary the config is provided via the `Any` type in the`typed_config` field. 

See comments in [any.proto](https://github.com/protocolbuffers/protobuf/blob/d4c5992352aae1ed18f44c1a40d2149006bf8704/src/google/protobuf/any.proto#L94-L111) from the protobuf project for additional documentation.  

```yaml
{{ simpleProtoYAML "clutch.config.gateway.v1.Service" }}
```

Example with `clutch.service.authn` and environment variables.
```yaml title="clutch-config.yaml"
...
services:
  - name: clutch.service.authn
    typed_config:
      "@type": types.google.com/clutch.config.service.authn.v1.Config
      oidc:
        issuer: ${OIDC_ISSUER}
        client_id: ${OIDC_CLIENT_ID}
        client_secret: ${OIDC_CLIENT_SECRET}
        redirect_url: "http://localhost:8080/v1/authn/callback"
        session_secret: ${AUTHN_SESSION_SECRET}
...
```

:::tip
For now, docs for each component's configuration are not auto-generated. In order to determine the configuration specification for a component, check at the well-known path in [api/config/](https://github.com/lyft/clutch/blob/main/api/config)

In the example above, the configuration schema is specified in [api/config/service/authn/v1/authn.proto](https://github.com/lyft/clutch/blob/main/api/config/service/authn/v1/authn.proto).
:::

## Frontend

The Clutch frontend requires configuration at build time to determine which installed workflows to register and allows for users to override default values.

A custom gateway generated from the [scaffolding tool](/docs/development/custom-gateway) will have a `register-workflows` script target in `frontend/package.json`. This script calls out to `@clutch-sh/tools` to parse the custom gateway's config file and register the found workflows. It expects the frontend config file at the path `frontend/src/clutch.config.js`.

Example:

```jsx title="frontend/src/clutch.config.js"
module.exports = {
  "@clutch-sh/ec2": {
    terminateInstance: {
      trending: true,
      componentProps: {
        resolverType: "clutch.aws.ec2.v1.Instance",
      },
    },
    resizeAutoscalingGroup: {
      componentProps: {
        resolverType: "clutch.aws.ec2.v1.AutoscalingGroup",
      },
    },
  },
  "@lyft/private-workflow": {
    example: {},
  },
};
```

In the configuration above there are some open source workflows registered, in this case with overrides for their `trending` values. Notice how these workflows also have a `componentProps` field specified. Some workflows will require prop values that are specific to the user. Without them the workflow will not register on the app even if listed in the config file. Take the `@clutch-sh/ec2` package as an example; both the `terminateInstance` and `resizeAutoscalingGroup` workflows require a `resolverType` prop.

If a config is invalid a warning will be emitted in the console denoting which workflow is misconfigured along with the which required props are missing, for example:

```
[@clutch-sh/ec2][instance/terminate] Not registered:
  Invalid config - missing required component props resolverType
```

It's important to note that only packages which are installed will be included, even if they are listed in the config file.
