ALTER TABLE experiment_config DROP CONSTRAINT experiment_config_pkey;
ALTER TABLE experiment_config ALTER COLUMN "id" SET DATA TYPE BIGINT USING id::BIGINT;
ALTER TABLE experiment_run DROP CONSTRAINT experiment_run_pkey;
ALTER TABLE experiment_run ALTER COLUMN "id" SET DATA TYPE BIGINT USING id::BIGINT;
ALTER TABLE experiment_run ALTER COLUMN "experiment_config_id" SET DATA TYPE BIGINT USING experiment_config_id::BIGINT;
