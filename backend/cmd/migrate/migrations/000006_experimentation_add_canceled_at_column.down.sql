ALTER TABLE experiment_run DROP COLUMN cancellation_time;
ALTER TABLE experiment_run ADD COLUMN scheduled_end_time TIMESTAMP WITH TIME ZONE;