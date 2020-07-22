---
title: Authentication & Authorization
{{ .EditURL }}
---

Clutch has modular support for authentication (i.e. `authn`) and authorization (i.e. `authz`). This will allow you to adapt these primitives to your environment without completely rewriting the core logic.

Currently Clutch supports Open ID Connect (OIDC) for authentication and ships with an RBAC engine for authorization.

### Authentication

The authentication components work together to create a valid session for a user from a third-party authentication provider.

The session information is stored in a JSON Web Token (JWT) for use by the authorization components, auditing components, etc.

#### Components

| Name | Description |
| --- | --- |
| `clutch.service.authn` | Handles token exchange with the authentication provider to verify the users identity and signs Clutch's own JWT. |
| `clutch.module.authn` | Provides the `callback` and `login` endpoints for the token exchange, which in turn call the authn service. |
| `clutch.middleware.authn` | Validates the JWT on incoming requests and inserts authentication information such as the user ID into the request context.  |

#### Configuration

In order to use Clutch's native authentication, **registration of all three `authn` components is required** in the gateway configuration.

```yaml title="clutch-config.yaml"
gateway:
  ...
  middleware:
    ...
    // highlight-next-line
    - name: clutch.middleware.authn
modules:
  ...
  // highlight-next-line
  - name: clutch.module.authn
services:
  ...
  // highlight-start
  - name: clutch.service.authn
    typed_config:
      "@type": types.google.com/clutch.config.service.authn.v1.Config
      oidc:
        issuer: https://provider.example.com
        client_id: ${CREDENTIALS_OIDC_CLIENT_ID}
        client_secret: ${CREDENTIALS_OIDC_CLIENT_SECRET}
        redirect_url: "${BASE_URL}/v1/authn/callback"
      session_secret: ${CREDENTIALS_SESSION_SECRET}
  // highlight-end
```

#### Customization

Clutch currently supports OIDC, for example with Okta. It is possible to support other authentication providers by extending or swapping out the authn service.

Furthermore, Clutch has support for a `groups` field in the claims. By default, no groups are appended to the claim. The authn service's `Provider` interface has a `WithClaimsFromOIDCTokenFunc` function that can be used to override claims derivation from the provider. At Lyft, we use this to call an internal service that provides the group IDs for a user.

### Authorization

The authorization components work together to provide an RBAC engine for Clutch.

#### Components

| Name | Description |
| --- | --- |
| `clutch.service.authz` | Evaluates policies to determine whether to allow or deny an action. |
| `clutch.middleware.authz` | Calls the authz service on each request to determine whether to allow or deny an action. |
| `clutch.module.authz` | Provides an endpoint for testing an action to see in advance whether it will be allowed or denied. Not currently required. |

#### Configuration

```yaml title="clutch-config.yaml"
gateway:
  middleware:
    ...
    - name: clutch.middleware.authn
    // highlight-start
    # note: authz must come after authn so user info is available
    - name: clutch.middleware.authz
    // highlight-end
services:
  ...
  // highlight-start
  - name: clutch.service.authz
    typed_config:
      "@type": types.google.com/clutch.config.service.authz.v1.Config
      role_bindings:
        - to: [superuser]
          principals:
            - user: alice@example.com
      roles:
        - role_name: superuser
          policies:
            - policy_name: allow-all
              method: "*"
  // highlight-end
```

The RBAC engine allows for resource-level rules based on the [API annotations](/docs/advanced/security-auditing#api-annotations). For more details on the configuration see the protobuf defintion for [clutch.config.service.authz.v1.Config](https://github.com/lyft/clutch/blob/main/api/config/service/authz/v1/authz.proto).

#### Customization

The built-in authz service, which provides a simple RBAC engine, can be swapped out with any other type of authorization framework.