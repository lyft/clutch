alter table experiment_run alter column "id" SET DATA type varchar(100) using id::varchar;
alter table experiment_run alter column "experiment_config_id" SET DATA type varchar(100) using experiment_config_id::varchar;
ALTER TABLE experiment_run ADD PRIMARY KEY (id);
alter table experiment_config alter column "id" SET DATA type varchar(100) using id::varchar;
ALTER TABLE experiment_config ADD PRIMARY KEY (id);
