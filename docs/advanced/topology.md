---
title: Infrastructure Topology Features
{{ .EditURL }}
---

The Topology feature set enables enables core Clutch capabilities,
as well as enabling API's that can be used based on your needs.

One of the main goals of the toplogy service was to create a caching mechanism,
that would aggregate and store all aspects of infratructure with the ability
for it to be easily extended.

This allows implmenations to write feature that are not bound to service providers API ratelimits, lantancies and etc.
Giving us the control to provide a consistent user experience as the userbase scales.

### Topology API's


### Topology Caching


### How to extend the Topology Cache


### Autocomplete




I wanted to give a more detailed update on this issue. This feature is now complete and the resolvers for AWS and Kubernetes have the autocomplete functionality for all the supported resource at this time.

However I need to write some documentation on how to utilize it. I'll briefly describe here how you can set this up and a little more detail about the architecture so you could utilize these feature in a verity of different ways.

Autocomplete is powered by the [Topology service](https://github.com/lyft/clutch/tree/main/backend/service/topology), the Topology service has a paginated [Search API](https://github.com/lyft/clutch/blob/main/api/topology/v1/topology_api.proto#L26-L32)  which give us the ability to search the `topology_cache` [table](https://github.com/lyft/clutch/blob/main/backend/cmd/migrate/migrations/000005_create_topology_cache_table.up.sql).

In order to populate this table with resource (aws instance, k8s pods etc) from your infrastructure I have built what is known as the Topology Cache, to enabled the cache you can provide the [cache configuration](https://github.com/lyft/clutch/blob/main/api/config/service/topology/v1/topology.proto#L10-L18) for topology service like so.

```
  - name: clutch.service.topology
    typed_config:
      "@type": types.google.com/clutch.config.service.topology.v1.Config
      cache: {}
```

It's important to load / specify the `clutch.service.topology` last in your clutch configuration, this is because the topology service will iterate through the service registry to see if a service satisfies the `CacheableTopology` interface.

https://github.com/lyft/clutch/blob/49f42bbf300faec5279b21d340ac80f8f3e00962/backend/service/topology/topology.go#L45-L56

If satisfied, the topology cache will start ingesting resources provided by the service and storing it in our `topology_cache` table to utilize for features like the [Search API](https://github.com/lyft/clutch/blob/main/api/topology/v1/topology_api.proto#L26-L32) which powers autocomplete.

For completeness we can look at the Kubernetes service [caching implantation](https://github.com/lyft/clutch/blob/main/backend/service/k8s/cache.go#L28-L49) as an example of how this interface is satisfied.

Now to recap, how do you setup autocomplete? 
1) Enable the postgres service `clutch.service.db.postgres`
2) Enable the topology service `clutch.service.topology` like so.
```
  - name: clutch.service.topology
    typed_config:
      "@type": types.google.com/clutch.config.service.topology.v1.Config
      cache: {}
```
3) Thats it!


Im going to create another issue to track the documentation of all these new features and systems as there was a lot built that could be leveraged in a verity of different ways.

Let me know if you have any questions!
