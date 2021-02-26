ALTER TABLE topology_cache DROP CONSTRAINT topology_cache_pkey;
ALTER TABLE topology_cache ADD CONSTRAINT topology_cache_id_key UNIQUE (id);
