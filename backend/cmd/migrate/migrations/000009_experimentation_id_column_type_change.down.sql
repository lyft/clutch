ALTER TABLE experiment_config DROP CONSTRAINT experiment_config_pkey;
alter table experiment_config alter column "id" SET DATA type BIGINT using id::BIGINT;
ALTER TABLE experiment_run DROP CONSTRAINT experiment_run_pkey;
alter table experiment_run alter column "id" SET DATA type BIGINT using id::BIGINT;
alter table experiment_run alter column "experiment_config_id" SET DATA type BIGINT using experiment_config_id::BIGINT;
