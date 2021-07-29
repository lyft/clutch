---
title: Topology API's
{{ .EditURL }}
---

The Topology feature set powers core Clutch capabilities such as autocomplete,
as well as providing API's that can be leveraged for a multitude of purposes.

One of the main goals of the topology service was to create an extensable caching mechanism,
to store all aspects of infrastructure with the ability for it to be easily accessed via APIs.
This allows implementers to write features that are not bound by service provider API rate limits and latencies and allows for Clutch to provide a consistent user experience that scales with your users.
Giving Clutch the control to provide a consistent user experience that you can control as your user base scales.

## Topology Caching

At the core of the topology service is its caching functionality which is the foundation that powers its APIs.

You can enable the cache by adding the following to your clutch configuration.

```yaml title="clutch-config.yaml"
services:
  ...
  # The topology services requires the postgres datastore to be configured
  - name: clutch.service.topology
    typed_config:
      "@type": types.google.com/clutch.config.service.topology.v1.Config
      // highlight-next-line
      cache: {}
```
There are additional [configuration](https://github.com/lyft/clutch/blob/c3097e5ad477952bb4bb90cc1fb5a126d7434565/api/config/service/topology/v1/topology.proto#L14-L28) options that can be tuned if necessary.

### Leader Election

Clutch elects a leader to handle the caching operations to ensure that only one instance
performs write heavy operations.

## How to extend the Topology Cache

Extending the topology cache for private gateways can be done by implementing the `CacheableTopology` [interface](https://github.com/lyft/clutch/blob/c3097e5ad477952bb4bb90cc1fb5a126d7434565/backend/service/topology/topology.go#L46-L57).

```go
type CacheableTopology interface {
  CacheEnabled() bool
  StartTopologyCaching(ctx context.Context, ttl time.Duration) (<-chan *topologyv1.UpdateCacheRequest, error)
}
```

If implemented, the topology service will begin caching resources provided by the services which implement the interface, storing them in the `topology_cache` table to utilize for features such as the [Search API](https://github.com/lyft/clutch/blob/c3097e5ad477952bb4bb90cc1fb5a126d7434565/api/topology/v1/topology_api.proto#L26-L32) which powers [autocomplete](#autocomplete).

An example of a the Topology interface being implemented can be found in he Kubernetes service [caching implementation](https://github.com/lyft/clutch/blob/c3097e5ad477952bb4bb90cc1fb5a126d7434565/backend/service/k8s/cache.go#L28-L61).

:::caution Config Load Order
It's important to load the `clutch.service.topology` last or near the bottom in your Clutch configuration. This is required since the topology service will iterate through the service registry at load time to see if any services implement the `CacheableTopology` interface to determine if caching should be enabled.
:::

## Autocomplete

Once the topology service and module are configured autocomplete for all known Clutch resources and workflows will be enabled by default.

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

If you have custom types within your gateway you can enable autocomplete by following the steps below.

1. Enable the searchable annotation on the resolver proto, the [example](https://github.com/lyft/clutch/blob/540f0acfb4809acb938e0fc8f52debf2868c9b1c/api/resolver/k8s/v1/k8s.proto#L11-L15) below is for Kubernetes pods.

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

2. You will likely have a custom resolver for custom types to satisfy the autocomplete function that is present on the resolver interface.
There are many examples throughout Clutch you can reference including [the Kubernetes resolver implementation](https://github.com/lyft/clutch/blob/main/backend/resolver/k8s/k8s.go#L247-L273).
You will notice that the autocomplete implementation for the resolver is almost identical to the AWS and Core Clutch resolver implementations.

3. Enabling autocomplete on the frontend varies depending on the implementation of the workflow.
if you are using the resolver component, autocomplete will be enabled without any additional effort.
However if you have a more custom workflow or page, you enable autocomplete on any `TextField` by adding the following.

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
