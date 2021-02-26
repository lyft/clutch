ALTER TABLE topology_cache DROP CONSTRAINT topology_cache_id_key;
ALTER TABLE topology_cache ADD PRIMARY KEY (id, resolver_type_url);
