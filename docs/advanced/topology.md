---
title: Topology API's
{{ .EditURL }}
---

The Topology feature set powers core Clutch capabilities such as autocomplete,
as well as providing API's that can be leveraged for a multitude of purposes.

One of the main goals of the topology service was to create an extensable caching mechanism,
to store all aspects of infratructure with the ability for it to be easily accessed via API's.
This allows implementors to write features that are not bound by service provider API rate limits and latencies and allows for Clutch to provide a consistent user experience that scales with your users.
Giving Clutch the control to provide a consistent user experience that you can control as your user base scales.

## Topology Caching

At the core of the topology service is its caching functionality which is the foundation that powers its API's.

You can enable the cache by adding the following to your clutch configuration.

```yaml title="clutch-config.yaml"
services:
  ...
  # The topology services does require the postgres datastore to be configured
  - name: clutch.service.topology
    typed_config:
      "@type": types.google.com/clutch.config.service.topology.v1.Config
      // highlight-next-line
      cache: {}
```

### Leader Election

Currently Clutch elects a leader to handle the caching operations so that,
only one instances performs write heavy operations.

## How to extend the Topology Cache

Extending the topology cache for private gateways can be done by satisfying the `CacheableTopology` [interface](https://github.com/lyft/clutch/blob/c3097e5ad477952bb4bb90cc1fb5a126d7434565/backend/service/topology/topology.go#L46-L57).

```go
type CacheableTopology interface {
  CacheEnabled() bool
  StartTopologyCaching(ctx context.Context, ttl time.Duration) (<-chan *topologyv1.UpdateCacheRequest, error)
}
```

If satisfied, the topology service will start caching resources provided by the services satisfiying the interface and storing it in the `topology_cache` table to utilize for features like the [Search API](https://github.com/lyft/clutch/blob/c3097e5ad477952bb4bb90cc1fb5a126d7434565/api/topology/v1/topology_api.proto#L26-L32) which powers [autocomplete](#autocomplete).

For completeness we can look at the Kubernetes service [caching implementation](https://github.com/lyft/clutch/blob/c3097e5ad477952bb4bb90cc1fb5a126d7434565/backend/service/k8s/cache.go#L28-L61) as an example of how this interface is satisfied.

:::caution Config Load Order
It's important to load the `clutch.service.topology` last or near the bottom in your Clutch configuration, this is because the topology service will iterate through the service registry to see if any satisfies the `CacheableTopology` interface to enable caching for those services.
:::

## Autocomplete

Once the topology service and module are configured, autocomplete for all known Clutch resources and workflows will be enabled by default.

```yaml title="clutch-config.yaml"
modules:
  // highlight-next-line
 - name: clutch.module.topology
services:
  ...
  # The topology services does require the postgres datastore to be configured
  // highlight-start
  - name: clutch.service.topology
    typed_config:
      "@type": types.google.com/clutch.config.service.topology.v1.Config
      cache: {}
  // highlight-end
```


### Autocomplete for custom resolver types

Enabling autocomplete for custom types you have defined in your gateway can be achieved by doing the following.

1.Enable the searchable annotation on the resolver proto, the [example](https://github.com/lyft/clutch/blob/540f0acfb4809acb938e0fc8f52debf2868c9b1c/api/resolver/k8s/v1/k8s.proto#L11-L15) below is for Kubernetes pods.

```protobuf
message PodID {
  option (clutch.resolver.v1.schema) = {
    display_name : "pod ID"
    // highlight-next-line
    search : {enabled : true}
  };
  ...
}
```

2.You will likely have a custom resolver for custom types, for this will have to satisfy the autocomplete function that is present on the resolver interface.
There are many examples in Clutch you can reference, [here](https://github.com/lyft/clutch/blob/main/backend/resolver/k8s/k8s.go#L247-L273) is the Kubernetes resolver implementation.
You will notice that the autocomplete function implementation for the resolver are just about identical for AWS and Core Clutch resolvers.

3.Enabling autocomplete on the frontend varies depending on the implementation of the workflow.
if you are using the wizard, autocomplete will be enabled without any additional effort.
However if you have a more custom workflow or page, you can do the following to enable autocomplete on any `TextField`.

```typescript
const autoComplete = async (search: string): Promise<any> => {
  // Check the length of the search query as the user might empty out the search
  // which will still trigger the on change handler
  if (search.length === 0) {
    return { results: [] };
  }

  const response = await client.post("/v1/resolver/autocomplete", {
    // Replace with the type you want to autocomplete on
    // highlight-next-line
    want: `type.googleapis.com/clutch.core.project.v1.Project`,
    search,
  });

  return { results: response?.data?.results || [] };
};


<TextField
  name="sometextfield"
  // Specify the autocompleteCallback prop for the TextField component.
  // highlight-next-line
  autocompleteCallback={v => autoComplete(v)}
/>
```
