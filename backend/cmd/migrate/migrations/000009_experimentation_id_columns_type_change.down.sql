BEGIN;
LOCK TABLE experiment_config, experiment_run IN ACCESS EXCLUSIVE MODE;
ALTER TABLE experiment_run ALTER COLUMN "experiment_config_id" SET DATA TYPE BIGINT USING experiment_config_id::BIGINT;
ALTER TABLE experiment_run ALTER COLUMN "id" SET DATA TYPE BIGINT USING id::BIGINT;
ALTER TABLE experiment_config ALTER COLUMN "id" SET DATA TYPE BIGINT USING id::BIGINT;
COMMIT;
