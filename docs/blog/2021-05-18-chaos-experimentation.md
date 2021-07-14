---
title: Chaos Experimentation
authors:
  - name: Kathan Shah
    url: https://github.com/kathan24
    avatar: https://avatars.githubusercontent.com/u/5263542?v=4
    twitter_username: kathan24
description: Chaos Experimentation, an open-source framework built on top of Envoy Proxy.
image: https://miro.medium.com/max/1400/1*Xi46XIWByMV7PUePnUIZsg.png
hide_table_of_contents: false
---

import useBaseUrl from '@docusaurus/useBaseUrl';

Services are bound to degrade. It’s a matter of when, not if. In a distributed system where there are many interdependent microservices, it is increasingly difficult to know what will happen when a service is unavailable, latency goes up, or when the success rate drops. Usually, companies find out the hard way when it happens in production and it affects their customers. This is where [Chaos Engineering](https://principlesofchaos.org/) helps us.

<!--truncate-->

---

**Note**: *[This article](https://eng.lyft.com/chaos-experimentation-an-open-source-framework-built-on-top-of-envoy-proxy-df87519ed681) was originally published at [eng.lyft.com](https://eng.lyft.com/).*

---

***Chaos Engineering is the discipline of experimenting on a system in order to build confidence in the system’s capability to withstand turbulent conditions in production.***

By regularly experimenting with service degradation in a controlled production environment, we can preemptively validate our system's resiliency and uncover issues. As Lyft grew, we quickly realized the importance of Chaos Engineering. Since all the communications between Lyft services run through [Envoy Proxy](https://www.envoyproxy.io/docs/envoy/latest/), it seemed like a great choice to leverage it for running chaos experiments like fault injection.

## A couple of years ago at Lyft

Lyft previously performed fault injection experiments using [Envoy’s runtime](https://www.envoyproxy.io/docs/envoy/latest/configuration/operations/runtime#runtime) (disk layer). Engineers ran a CLI command that generated runtime files locally. Once runtime files were committed to GitHub, they got deployed by writing to the local file system of a cluster of hosts. Envoy read these files into memory and injected faults into the requests. Engineers had to repeat the same process when they wanted to terminate the fault. This worked well when Lyft was at an early stage and performed one-off experiments to prepare for a surge in traffic on days like Halloween, New Year’s Eve, etc. However, this had its drawbacks when running at scale.

- **High touch**

    Many times, engineers needed to perform multiple steps in their experiments, such as injecting faults for 1% of requests and then slowly increasing to ensure safety. This involved multiple commits to GitHub, and this process was very cumbersome.

- **Long time for faults to be injected**

    Once the runtime changes were merged, runtime deployment took a few minutes to complete. This could be risky when the experiment caused a real production issue and the engineer wanted to terminate the experiment right away.

- **Poor insight into all active faults**

    It was very difficult to find out how many faults were running at any given time. This info is very important to know when looking at all experiments from a bird’s-eye view.

- **External dependency — GitHub**

    This process had a dependency on GitHub. If GitHub were to go down (which could happen since it’s a service, after all), the active faults would persist in the system for a long, undefined period of time. This was a big drawback and could jeopardize our business.

## Launching the Chaos Experimentation Framework

The Chaos Experimentation Framework (CEF) is an open-source framework built on top of [Envoy Proxy](https://www.envoyproxy.io/docs/envoy/latest/) in Clutch. Clutch made perfect sense for integrating the CEF since we could leverage built-in features like rich and responsive UI, role-based access controls, audit logging, database layers, stats sinks, etc.

The CEF provides a powerful backend and user-friendly UI which gives high confidence to engineers when they perform experiments. The backend injects faults within a few milliseconds after starting the experiment from the UI. This is achieved using Envoy Proxy’s xDS APIs to transmit fault configuration rather than relying on deploying configuration on disk to every machine. It was designed to quickly inject faults without any dependencies.

![Server Chaos Experiment Flow](https://miro.medium.com/max/1400/1*ArsUFuV7EKSV1HcEZKAuew.gif)

## Benefits

With the CEF, we have already seen several benefits at Lyft:

- **Self-serve**

    This framework is fully self-service, and Lyft engineers can perform fault injection experiments on their services quickly with the click of a button.

- **Integration with CI/CD pipeline**

    We run experiments on all deployments to ensure that deploying new code will not affect a service’s resilience to failures. This ensures that resiliency is validated regularly across the system.

- **Ensure client resiliency**

    Although we work to ensure the resiliency of our service mesh, at Lyft it’s equally important to make sure our mobile clients have a fallback plan when things don’t work as expected. Usually, product flows on mobile clients are tested under ideal conditions, i.e., the “happy path”. However, flows in non-ideal paths are hard to test in the QA environment. By regularly running fault injection experiments on our mobile client endpoints, we can minimize the chance that our customers are affected when a backend service is degraded.

- **Faster service tier auditing**

    Lyft services are categorized in different tiers (Tier 0 — Tier 3) based on how business-critical each service is. With Tier 0 being the most important, we inject faults in Tier 1 services that have Tier 0 downstream dependencies. With this method, we are simulating the situation where one specific Tier 1 service is experiencing a degradation, and observing how this degradation affects that service’s Tier 0 downstream dependents. Ideally, there should be no hard dependency on a Tier 1 service from Tier 0. However, we discovered some situations where this was not the case. With the CEF, this tier auditing is faster than ever.

- **Validating fallbacks to third party external services**

    Usually, it’s very hard to test failures from external services (like Mapbox, Google Maps, DynamoDB, Redis, etc) until they happen. With the CEF, Lyft engineers have been proactively validating fallback logic in case external services become degraded.

- **Keeping observability and configuration up-to-date**

    By running periodic fault injection tests, engineers preemptively tune service alarms, timeouts, and retry policies. Along with these tune-ups, they also make sure their services’ stats and logging provide clear indications of root causes when issues arise.

*“The Chaos Experimentation framework is very simple and straightforward to use. It allows service owners to determine the resilience of their microservices.”* — Testimonial from a Lyft Engineer

## Architecture

![Architecture of Chaos Experimentation Framework and its interaction with Envoy service mesh](https://miro.medium.com/max/1400/1*Xi46XIWByMV7PUePnUIZsg.png)

There are two major components in the framework — the backend server and the [Envoy xDS management server](https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol).

The backend server is responsible for all the CRUD operations of the experiments. It stores experiments in the tables of its Postgres database.

The other component that the framework ships with is an xDS management server. The management server consists of [Extension Configuration Discovery Service](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration#ecds) (ECDS) and [Runtime Discovery Service](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/operations/dynamic_configuration#rtds) (RTDS) APIs of Envoy Proxy. Either of the xDS APIs can be used to perform fault injection experiments: With RTDS, one can make changes to runtime-specific faults. Conversely, ECDS allows for changes to the entire fault filter to perform any custom experiments. When Envoy in the mesh boots up, it creates a bi-directional gRPC stream with the management server. Below is a code snippet for an RTDS config in Envoy:

```yaml
layered_runtime:
  layers:
  - name : rtds
    rtds_layer:
    name: rtds
    rtds_config:
      resource_api_version: v3
      api_config_source:
        api_type: GRPC
        transport_api_version: v3
        grpc_services:
          envoy_grpc: {cluster_name: clutchxds}
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

The management server polls Postgres at a regular cadence to get all of the active experiments. It then forms a runtime resource with a TTL that is sent to its respective clusters. Hence, the propagation of faults only takes a few milliseconds.

The framework itself is very resilient. It automatically terminates experiments when the success rate of service drops beyond a configured threshold. Additionally, in the case where the management server itself becomes degraded, all experiments are automatically disabled without any intervention from engineers.

The entire framework, like Clutch and Envoy Proxy, is config-driven. One can choose to use ECDS or RTDS, tune the polling duration to Postgres, provide runtime prefixes, tune resource TTL times, etc.

## Road ahead

There is still a lot of work to be done in this space to prevent system degradation from affecting customers. Here are some of our ideas for future improvements:

- **[Scheduling of experiments](https://github.com/lyft/clutch/issues/1356)**

    Scheduling would allow us to run experiments 24/7 and provide more confidence that the system is up-to-date and resilient.

- **[Real-time stats](https://github.com/lyft/clutch/issues/1357) with [Envoy’s Load Reporting Service API](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/load_stats/v3/lrs.proto)**

    To perform more aggressive experiments, there needs to be a tight metric-driven feedback system to ensure that we can quickly terminate experiments before they affect our users.

- **[Squeeze experiments](https://github.com/lyft/clutch/issues/1369)**

    Squeeze tests would allow us to route additional traffic to a particular host in a given service and help determine the maximum number of concurrent requests that can be served by that host in a cluster. Based on squeeze experiments, engineers can set the scaling threshold of the service and its [circuit-breaking threshold](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/circuit_breaking).

- **[Enhanced Redis support](https://github.com/lyft/clutch/issues/1368)**
 
    The framework already provides basic support for injecting Redis faults. However, more advanced features would allow for injecting faults based on a certain set of Redis commands and performing latency injection or connection failure experiments.

If you’re ready to get started or contribute, check out these resources:
- [Documentation](https://clutch.sh/docs/advanced/chaos-experimentation)
- [Code in Clutch](https://github.com/lyft/clutch)
- Join Lyft — [apply here](https://www.lyft.com/careers)

## Thank you
Chaos Experimentation, an open-sourced framework, would not have been possible without the contributions and hard work from many engineers at Lyft including Alexander Herbert, Ansu Kar, Bill Gallagher, Daniel Hochman, Derek Schaller, Don Yu, Gastón Kleiman, Ivan Han, Jingwei Hao, Jyoti Mahapatra, Miguel Juárez, Mike Cutalo, Rafal Augustyniak, Snow Pettersen, and Vijay Rajput. Special thanks to Patrick Sunday, Martin Conte Mac Donell, Matt Klein, Polly Peterson, Michael Rebello, and Pete Morelli.
