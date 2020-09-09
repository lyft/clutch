CREATE TABLE topology_cache (
  -- id: This is the resource identifier, could be a pod-id, aws instance-id, this value must be unique
  id varchar UNIQUE,

  -- data: json blob of resource
  data JSONB,

  -- resolver_type_url: is the resolver proto type for example: `type.googleapis.com/clutch.k8s.v1.Deployment`
  resolver_type_url varchar,

  -- metadata: Give us the ability to define and query different dimensions of our cache
  -- being type JSONB this the structure will be dependent based on usecase.
  metadata JSONB,

  -- updated_at: Keeps track of when the cache entry was last created / update so we can expire old entries
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX topology_cache_id_search_idx ON topology_cache (id varchar_pattern_ops);
CREATE INDEX topology_cache_metadata_idx ON topology_cache USING GIN (metadata jsonb_path_ops);
