CREATE OR REPLACE FUNCTION set_timestamp() RETURNS TRIGGER LANGUAGE plpgsql AS '
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
';

CREATE FUNCTION clean_cache() RETURNS TRIGGER LANGUAGE plpgsql AS '
BEGIN
  DELETE FROM topology_cache_data WHERE updated_at <= NOW() - INTERVAL ''60m'';
  RETURN NULL;
END;
';

-- Tables
CREATE TABLE topology_cache_data (
  -- id: This is the resource identifier, could be a pod-id, aws instance-id, this value must be unique
  id varchar UNIQUE,
  
  -- id_text_search: This is the same value as the id but in tsquery type for better indexing and fuzzy lookups
  id_text_search tsquery,
  
  -- data: json blob of resource
  data JSONB,
  
  -- resolver_type: is the infrastructure type that the resource belongs to, this could be `kubernetes` or `aws`
  resolver_type varchar,
  
  -- metadata: Give us the ability to define and query different dimensions of our cache
  -- being type JSONB this the structure will be dependent based on usecase.
  metadata JSONB,
  
  -- updated_at: Keeps track of when the cache entry was last created / update so we can expire old entries
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX topology_cache_data_labels_idx ON topology_cache_data USING GIN (metadata jsonb_path_ops);
CREATE INDEX topology_cache_data_text_search_idx ON topology_cache_data USING GIST (id_text_search);

CREATE TRIGGER trigger_set_timestamp
BEFORE UPDATE ON topology_cache_data
FOR EACH ROW 
EXECUTE PROCEDURE set_timestamp();

CREATE TRIGGER trigger_clear_cache
AFTER INSERT ON topology_cache_data
EXECUTE PROCEDURE clean_cache();

INSERT INTO topology_cache_data (id, id_text_search, data, metadata, resolver_type) VALUES (
    'clutch-main-123',
    'clutch-main-123',
    '{
    	 "spec": {"bits": 1}
     }',
    '{
        "app":"clutch",
        "environment":"staging",
        "facet":"main",
        "facet-type":"service"
     }',
    'kubernetes'
);

INSERT INTO topology_cache_data (id, id_text_search, data, metadata, resolver_type) VALUES (
    'clutch-main-1234',
    'clutch-main-1234',
    '{
    	 "spec": {"bits": 2}
     }',
    '{
        "app":"clutch",
        "environment":"staging",
        "facet":"main",
        "facet-type":"service"
     }',
    'kubernetes'
);

INSERT INTO topology_cache_data (id, id_text_search, data, metadata, resolver_type, updated_at) VALUES (
    'clutch-rtds-12346',
    'clutch-rtds-12346',
    '{
    	 "spec": {"bits": 2}
     }',
    '{
        "app":"clutch",
        "environment":"staging",
        "facet":"rtds",
        "facet-type":"service"
     }',
    'kubernetes',
	'2020-07-27T16:30:39.595Z'
);
