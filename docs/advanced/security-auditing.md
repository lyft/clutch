---
title: Security Auditing
{{ .EditURL }}
---

import useBaseUrl from '@docusaurus/useBaseUrl';

Saving each event performed with Clutch for security auditing is a top-level concern, and the processes to accomplish this fit into the normal component architecture.

<img alt="Clutch Component Architecture" src={useBaseUrl('img/docs/audit-architecture-diagram.png')} width="75%" />

Clutch has configurable middleware to save data on each incoming request, a module to query for previous events, and services to store and forward events to new sinks. What information it stores for the request comes from API annotations on the request itself. The below table has the components ordered as they are in a call-chain.

| Component Name                | Description |
| ----------------------------- | ----------- |
| `clutch.module.audit`         | Module to retrieve audit events, called from outside of Clutch and calls the audit service. |
| `clutch.middleware.audit`     | Middleware that runs on each request. Calls the audit service to store request information. |
| `clutch.service.audit`        | Handler in front of storage services. Called to save or retrieve request information. |
| `clutch.service.audit.sink.*` | Calls the audit service to retrieve request information. Call out of Clutch to forward the events. |

The annotations are read by middleware for each request and response, which forwards the information to the audit service for storage. The audit service is responsible for persistence as well as forwarding to other sinks to fan-out the messages. Each step of saving and forwarding is configurable with filters based on event attributes.

### API Annotations

Clutch's auditing middleware reads the value of two API annotations: the request "type" to determine what sort of action is being performed on a resource, and the resource name to identify what item is being modified.

These annotations let customizable data on each request be logged and reviewed.

#### Action Types

Clutch will save the the "type" of action that was performed (create, read, update, or delete or CRUD operations) for each request. The request type comes from an API annotation on each RPC.

From the [`AmiiboAPI` example proto](/docs/development/feature#2-api-define-the-schema), to show that the `GetCharacters` endpoint is a "read" operation, an extra `option` annotation would be needed.

```protobuf title="api/amiibo/v1/amiibo.proto"
syntax = "proto3";

package clutch.amiibo.v1;

option go_package = "amiibov1";

import "google/api/annotations.proto";

// highlight-next-line
import "api/v1/annotations.proto";

service AmiiboAPI {
  rpc GetCharacters(GetCharactersRequest) returns (GetCharactersResponse) {
    option (google.api.http) = {
      post : "/v1/amiibo/getCharacters"
      body : "*"
    };

// highlight-next-line
    option (clutch.api.v1.action).type = READ;
  }
}
```

The annotation would look similar for the other supported types of changes.

#### Resource Names

Clutch will also resolve what resources were modified by each request to their unique identifiers. This is also accomplished through API annotations.

There are two different annotations used to mark the resources in a single request:
  - A `fields` annotation will indicate what fields contain resource names.
  - A `patterns` annotation will try template substitution with field values in order to create a name.

If we wanted to add these annotations to the Amiibo example, where the identifier for a character is a combination of its name and game series, we would make the following modifications to the API definition.

```protobuf title="api/amiibo/v1/amiibo.proto"
syntax = "proto3";

package clutch.amiibo.v1;

option go_package = "amiibov1";

import "google/api/annotations.proto";
import "validate/validate.proto";

import "api/v1/annotations.proto";

service AmiiboAPI {
  rpc GetCharacters(GetCharactersRequest) returns (GetCharactersResponse) {
    option (google.api.http) = {
      post : "/v1/amiibo/getCharacters"
      body : "*"
    };

    option (clutch.api.v1.action).type = READ;
  }
}

...

message GetCharactersResponse {
  // highlight-next-line
  option (clutch.api.v1.id).fields = "characters";

  repeated Character characters = 1;
}

message Character {
  // highlight-start
  option (clutch.api.v1.id).patterns = {
    type_url : "amiibo.v1.Character",
    pattern : "{name}/{game_series}"
  };
  // highlight-end

  string name = 1;
  string game_series = 2;
  string amiibo_series = 3;
  string image_url = 4;

  enum Type {
    UNSPECIFIED = 0;
    CARD = 1;
    FIGURE = 2;
    YARN = 3;
  }
  Type type = 5;
}
```

The first annotation (`fields = "characters"`) says that to get the set of resources contained in this message, to look at the `characters` field.

The second annotation (`patterns = ...`) says that the identifier for this message is built by combining the value of the `name` field with the value of the `game_series` field.

### Middleware

On each request, the [above annotations](#annotations) are read along with what endpoint was called, by whom, and when. When a response is being sent, any additional resources or final identifiers are merged with the initial list from the request.

All of this information is passed along to the audit [service](./###Service) to persist.

### Service

The audit service has two behaviors: write requests somewhere, and read them back out. It takes events from the middleware and saves them, and it also pushes them to later "sinks" for further processing.

Clutch ships with an implementation that uses Postgres as backing storage. With this, events that are saved can be filtered. At Lyft, we use this to prevent saving healthcheck information, since it is not interesting to us from a security perspective.

The configuration for this looks like:
```yaml title="backend/clutch-config.yaml"
...
services:
  ...
  // highlight-start
  - name: clutch.service.db.postgres
    typed_config:
      "@type": types.google.com/clutch.config.service.db.postgres.v1.Config
      connection:
        host: ${DB_HOST}
        port: ${DB_PORT}
        user: ${DB_USER}
        ssl_mode: DISABLE
        dbname: clutch
  - name: clutch.service.audit
    typed_config:
      "@type": types.google.com/clutch.config.service.audit.v1.Config
      db_provider: clutch.service.db.postgres
      sinks:
        - clutch.service.audit.sink.slack
  // highlight-end
```

### Sinks

Sinks asynchronously propagate events to other systems after they are persisted to Clutch's database.

Clutch ships with a logging sink as a scaffold for your own, as well as a sink for Slack.

Adding and customizing audit sinks lets you save or process infrastructure events however appropriate for your needs.

### Module

Clutch's audit events can also be viewed by querying the audit module if it is enabled.

### Example config

#### Local Adhoc Use

Clutch ships with an in-memory storage for events which allows it to be used without
setting up a Postgres database. It is not recommended to run this way outside of a trial
or adhoc temporary use. All history of actions taken with Clutch will be lost with the
process shutting down.


```yaml title="backend/clutch-config.yaml"
...
services:
  ...
  // highlight-start
  - name: clutch.service.audit.sink.logger
  - name: clutch.service.audit
    typed_config:
      "@type": types.google.com/clutch.config.service.audit.v1.Config
      local_auditing: true
      filter:
        denylist: true
        rules:
          - field: METHOD
            text: Healthcheck
      sinks:
        - clutch.service.audit.sink.slack
  // highlight-end
```



#### Recommended Production Setup

Below is sample configuration to show how the services described are enabled. Note that because services are instantiated in the order they are listed, order matters! Since the audit service depends on both the database and the sink, it needs to be listed after them.

```yaml title="backend/clutch-config.yaml"
...
services:
  ...
  // highlight-start
  - name: clutch.service.db.postgres
    typed_config:
      "@type": types.google.com/clutch.config.service.db.postgres.v1.Config
      connection:
        host: ${DB_HOST}
        port: ${DB_PORT}
        user: ${DB_USER}
        password: ${DB_PASSWORD}
        ssl_mode: DISABLE
        dbname: clutch
  - name: clutch.service.audit.sink.logger
  - name: clutch.service.audit
    typed_config:
      "@type": types.google.com/clutch.config.service.audit.v1.Config
      db_provider: clutch.service.db.postgres
      filter:
        denylist: true
        rules:
          - field: METHOD
            text: Healthcheck
      sinks:
        - clutch.service.audit.sink.slack
  // highlight-end
```
