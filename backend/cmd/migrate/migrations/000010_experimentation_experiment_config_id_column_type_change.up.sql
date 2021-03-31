BEGIN;
ALTER TABLE experiment_config ALTER COLUMN "id" SET DATA TYPE varchar(100) USING id::varchar;
ALTER TABLE experiment_config ADD PRIMARY KEY (id);
COMMIT;
