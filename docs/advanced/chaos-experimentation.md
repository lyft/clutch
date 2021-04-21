---
title: Chaos Experimentation Framework
{{ .EditURL }}
---

import useBaseUrl from '@docusaurus/useBaseUrl';

A framework to perform Chaos experiments built on top of [EnvoyProxy](https://www.envoyproxy.io/docs/envoy/latest/)

<img alt="Clutch Component Architecture" src={useBaseUrl('img/docs/chaos-experimentation.png')} width="75%" />

Chaos Experimentation Framework consists of few parts - frontend, a backend server and a xDS management server. 

The Frontend uses [Clutch's core frontend](https://clutch.sh/docs/development/frontend). It can be customized by using the frontend config. 

The Backend server is responsible for all the CRUD operations of the experiments. It stores experiments in its tables in the Postgres database. 

The xDS management server uses [go-control-plane](http://github.com/envoyproxy/go-control-plane) library and serves two Envoy APIs - [RTDS](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration#rtds) and [ECDS](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration#ecds). Either of these two APIs can be used to perform fault injection tests. With RTDS, you can make changes to runtimes specific to faults whereas with ECDS you can make changes to the entire fault filter to perform any custom experiments.

### Components

The framework needs below mentioned components to perform experiments starting from storing the data into the Postgres database for each incoming request all the way to pass the experiment values to the Envoys to inject faults.

| Component Name                                | Description |
| --------------------------------------------- | ----------- |
| `clutch.module.chaos.experimentation.api`     | Module that supports CRUD API for managing experiments like Create, Get, List, Cancel, etc|
| `clutch.module.chaos.serverexperimentation`   | Module which is responsible to processes the experiment data that it receives from database to be UI ready
| `clutch.module.chaos.experimentation.xds`     | Module which runs Envoy xDS management server which is responsible to propagate experiment values to Envoys
| `clutch.service.chaos.experimentation.store`  | Service that defines the data layer to perform all database operations for experiments  
| `clutch.service.db.postgres`                  | Service to connect to Postgres database

In order to use Chaos Experimentation Framework, **registration of all the above components is required**. 

It is recommended to run Envoy xDS management server (`clutch.module.chaos.experimentation.xds`) on a separate host.    

### Configuration

#### Frontend

The frontend of the framework is completely configurable. Below is a sample [frontend config](https://github.com/lyft/clutch/blob/6990b5aa8b1e6a47a33b28e2aaab9783e4e9d084/frontend/packages/app/src/clutch.config.js) which will show the list of experiments and as well as workflow to start/stop an experiment. 

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
  - name: clutch.module.chaos.lyftexperimentation
  - name: clutch.module.chaos.serverexperimentation
services:
  ...
  - name: clutch.service.db.postgres
      typed_config:
        "@type": types.google.com/clutch.config.service.db.postgres.v1.Config
        connection:
          host: ${RDS_HOST}
          port: ${RDS_PORT}
          user: ${RDS_USER}
          ssl_mode: REQUIRE
          dbname: ${RDS_NAME}
          password: ${CREDENTIALS_CLUTCH_TO_RDS_PASSWORD}
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
        host: ${RDS_HOST}
        port: ${RDS_PORT}
        user: ${RDS_USER}
        ssl_mode: REQUIRE
        dbname: ${RDS_NAME}
        password: ${CREDENTIALS_CLUTCH_TO_RDS_PASSWORD}
  - name: clutch.service.chaos.experimentation.store
```

Keep in mind that both backend config and xDS config needs to connect to the same Postgres database.

### Example Envoy config

When Envoy in the mesh boots up, it creates a bi-directional gRPC stream with the management server. Below is the sample Envoy configs for RTDS and ECDS which will initiate the connection to the xDS server 

#### RTDS 
```yaml title="envoy-rtds-config.yaml"
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
        
```

#### ECDS
```yaml title="envoy-ecds-config.yaml"
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
        '@type': type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault
        abort: 
          http_status: 503
          percentage: {denominator: HUNDRED, numerator: 0}
        abort_http_status_runtime: ecds_runtime_override_do_not_use.http.abort.http_status
        abort_percent_runtime: ecds_runtime_override_do_not_use.http.abort.abort_percent
        delay: 
          fixed_delay: 0.001s
          percentage: {denominator: HUNDRED, numerator: 0}
        delay_duration_runtime: ecds_runtime_override_do_not_use.http.delay.fixed_duration_ms
        delay_percent_runtime: ecds_runtime_override_do_not_use.http.delay.percentage
      apply_default_config_without_warming: false
      type_urls: 
        - type.googleapis.com/envoy.extensions.filters.http.fault.v3.HTTPFault
```

### Redis experiments 

To perform Redis experiments, there is a specific module that is used to process the Redis experiment data. You will need below component in addition to the above Experimentation components 

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