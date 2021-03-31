BEGIN;
LOCK TABLE experiment_run IN ACCESS EXCLUSIVE MODE;
ALTER TABLE experiment_run ALTER COLUMN "id" SET DATA TYPE varchar(100) USING id::varchar;
ALTER TABLE experiment_run ALTER COLUMN "experiment_config_id" SET DATA TYPE varchar(100) USING experiment_config_id::varchar;
ALTER TABLE experiment_run ADD PRIMARY KEY (id);
COMMIT;
