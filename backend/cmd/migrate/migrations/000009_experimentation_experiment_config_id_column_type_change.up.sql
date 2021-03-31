BEGIN;
LOCK TABLE experiment_config IN ACCESS EXCLUSIVE MODE;
ALTER TABLE experiment_config ALTER COLUMN "id" SET DATA TYPE varchar(100) USING id::varchar;
COMMIT;
