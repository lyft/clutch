ALTER TABLE experiment_config DROP CONSTRAINT experiment_config_pkey;
ALTER TABLE experiment_config ALTER COLUMN "id" SET DATA TYPE BIGINT USING id::BIGINT;
