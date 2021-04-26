---
title: Chaos Experimentation Framework
{{ .EditURL }}
---

import useBaseUrl from '@docusaurus/useBaseUrl';

A framework to perform Chaos experiments built on top of [EnvoyProxy](https://www.envoyproxy.io/docs/envoy/latest/)

<img alt="Clutch Component Architecture" src={useBaseUrl('img/docs/chaos-experimentation.png')} width="75%" />

Chaos Experimentation Framework consists of a few parts - frontend, a backend server and a xDS management server. 

The Frontend uses [Clutch's core frontend](https://clutch.sh/docs/development/frontend). It can be customized by using the frontend config. 

The Backend server is responsible for performing CRUD operations of the experimentation package - [CreateExperiment](https://github.com/lyft/clutch/blob/71f84e4bb3f642a17b831019f188a87dcc63f2cf/backend/module/chaos/experimentation/api/experimentation.go#L68), [GetExperiments](https://github.com/lyft/clutch/blob/71f84e4bb3f642a17b831019f188a87dcc63f2cf/backend/module/chaos/experimentation/api/experimentation.go#L134), [CancelExperimentRun](https://github.com/lyft/clutch/blob/71f84e4bb3f642a17b831019f188a87dcc63f2cf/backend/module/chaos/experimentation/api/experimentation.go#L123), etc. It stores experiments in its tables in the Postgres database. 

The xDS management server uses [go-control-plane](http://github.com/envoyproxy/go-control-plane) library and serves two Envoy APIs - [Runtime Discovery Service (RTDS)](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration#rtds) and [Extension Configuration Discovery Service (ECDS)](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration#ecds). Either of these two APIs can be used to perform fault injection tests. With RTDS, you can make changes to runtime specific to faults whereas with ECDS you can make changes to the entire fault filter to perform any custom experiments.

### Components

Below components are responsible to perform experiments starting from storing the data into the Postgres database for each incoming request all the way to passing the experiment values to the Envoys to inject faults.

| Component Name                                | Description |
| --------------------------------------------- | ----------- |
| `clutch.module.chaos.experimentation.api`     | Module that supports CRUD API for managing experiments like Create, Get, List, Cancel, etc|
| `clutch.module.chaos.serverexperimentation`   |  Module responsible for orchestrating server fault chaos experiments.
| `clutch.module.chaos.experimentation.xds`     | Module which runs Envoy xDS management server which is responsible for propagating chaos experiment configurations to Envoys
| `clutch.service.chaos.experimentation.store`  | Service that defines the data layer to perform all database operations for experiments  
| `clutch.service.db.postgres`                  | Service used to connect to Postgres database

In order to use Chaos Experimentation Framework, **registration of all the above components is required**. 

It is recommended to run Envoy xDS management server (`clutch.module.chaos.experimentation.xds`) on a separate host.    

### Configuration

#### Frontend

The frontend of the framework is completely configurable. Below is an [example frontend config](https://github.com/lyft/clutch/blob/6990b5aa8b1e6a47a33b28e2aaab9783e4e9d084/frontend/packages/app/src/clutch.config.js) which will show the list of experiments and as well as workflow to start/stop an experiment. 

```"@clutch-sh/experimentation": {
module.exports = {
  ...
  "@clutch-sh/experimentation": {
    listExperiments: {
      description: "Manage fault injection experiments.",
      trending: true,
      componentProps: {
        columns: [
          { id: "target", header: "Target" },
          { id: "fault_types", header: "Faults" },
          { id: "start_time", header: "Start Time", sortable: true },
          { id: "end_time", header: "End Time", sortable: true },
          { id: "run_creation_time", header: "Creation Time", sortable: true },
          { id: "status", header: "Status" },
        ],
        links: [
          {
            displayName: "Start Server Experiment",
            path: "/server-experimentation/start",
          },
        ],
      },
    },
    viewExperimentRun: {},
  },
  "@clutch-sh/server-experimentation": {
    startExperiment: {
      componentProps: {
        upstreamClusterTypeSelectionEnabled: true,
      },
      hideNav: true,
    },
  },
```

#### Backend Server

Below configuration will spin up all the required modules and services to store the data coming from frontend into Postgres database.

```yaml title="backend/clutch-config.yaml"
modules:
  ...
  - name: clutch.module.chaos.experimentation.api
  - name: clutch.module.chaos.serverexperimentation
services:
  ...
  - name: clutch.service.db.postgres
      typed_config:
        "@type": types.google.com/clutch.config.service.db.postgres.v1.Config
        connection:
          host: <RDS_HOST>
          port: <RDS_PORT>
          user: <RDS_USER>
          ssl_mode: REQUIRE
          dbname: <RDS_NAME>
          password: <RDS_PASSWORD>
  - name: clutch.service.chaos.experimentation.store
```

#### xDS Management Server 

Below is the configuration for spinning up xDS server. For details about the fields, take a look at the [xds config proto](https://github.com/lyft/clutch/blob/main/api/config/module/chaos/experimentation/xds/v1/xds.proto).

```yaml title="backend/clutch-xds-config.yaml"
...
modules:
  ...
  - name: clutch.module.chaos.experimentation.xds
    typed_config:
      "@type": types.google.com/clutch.config.module.chaos.experimentation.xds.v1.Config
      rtds_layer_name: <RTDS_LAYER_NAME>                     // "rtds_layer"
      cache_refresh_interval: <CACHE_REFRESH_INTERNAL>       // "5s"
      ingress_fault_runtime_prefix: <INGRESS_FAULT_PREFIX>   // "fault.http"
      egress_fault_runtime_prefix: <EGRESS_FAULT_PREFIX>     // "fault.http.egress"
      resource_ttl: <RESOURCE_TTL>                           // "20s"
      heartbeat_interval: <HEARTBEAT_INTERVAL>               // "5s"
      ecds_allow_list: <LIST_OF_ECDS_ENALBED_CLUSTERS>       // ["foo", "bar"]
services:
  - name: clutch.service.db.postgres
    typed_config:
      "@type": types.google.com/clutch.config.service.db.postgres.v1.Config
      connection:
        host: <RDS_HOST>
        port: <RDS_PORT>
        user: <RDS_USER>
        ssl_mode: REQUIRE
        dbname: <RDS_NAME>
        password: <RDS_PASSWORD>
  - name: clutch.service.chaos.experimentation.store
```

Keep in mind that both backend config and xDS config need to connect to the same Postgres database.

### Example Envoy config

When Envoy in the mesh boots up, it creates a bi-directional gRPC stream with the management server. Below is the sample Envoy configs for RTDS and ECDS which will initiate the connection to the xDS server. Checkout [Envoy Proxy docs](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/fault_filter) for details on Fault Injection

#### RTDS 
```yaml title="envoy.yaml"
...
layered_runtime:
  layers:
    - name: rtds
      rtds_layer:
        name: <RTDS_LAYER_NAME>
        rtds_config:
          api_config_source:
            api_type: GRPC
            grpc_services:
              envoy_grpc:
                cluster_name: <xDS_CLUSTER>
...
http_filters:
  - name: envoy.fault
    typed_config:
      "@type": "type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault"
      abort:
        percentage:
          numerator: 0
          denominator: HUNDRED
        http_status: 503
      delay:
        percentage:
          numerator: 0
          denominator: HUNDRED
        fixed_delay: 0.001s
...        
```

#### ECDS
```yaml title="envoy.yaml"
filters:
  ...
  http_filters:
    ...
    name: envoy.extension_config
    config_discovery:
      config_source:
        api_config_source: 
          api_type: GRPC
          grpc_services: 
            - envoy_grpc: 
                cluster_name: <xDS_CLUSTER>
          transport_api_version: V3
        initial_fetch_timeout: 10s
        resource_api_version: V3
      default_config: 
        "@type": "type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault"
          abort:
            percentage:
              numerator: 0
              denominator: HUNDRED
            http_status: 503
          delay:
            percentage:
              numerator: 0
              denominator: HUNDRED
            fixed_delay: 0.001s
      apply_default_config_without_warming: false
      type_urls: 
        - type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault
```

### Redis experiments 

To perform Redis experiments, there is a specific module that is used to process the Redis experiment data. You will need below component in addition to the above Experimentation components. Also, keep in mind that Redis experiments can be only performed with RTDS. 

| Component Name                                | Description |
| --------------------------------------------- | ----------- |
| `clutch.module.chaos.redisexperimentation`    | Module which is responsible to processes the data specifically for Redis experiments that it receives from database to be UI ready

#### Frontend Config

```"@clutch-sh/experimentation": {
module.exports = {
  "@clutch-sh/experimentation": {
    listExperiments: {
      ...
        links: [
          {
            displayName: "Start Redis Experiment",
            path: "/redis-experimentation/start",
          },
        ],
      },
    },
  },
  ...
  "@clutch-sh/redis-experimentation": {
    startExperiment: {
      hideNav: true,
    },
  },
```

#### Backend Config
```yaml title="clutch-config.yaml"
modules:
  ...
  - name: clutch.module.chaos.redisexperimentation
```
