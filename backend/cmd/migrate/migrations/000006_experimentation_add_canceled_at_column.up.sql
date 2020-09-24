ALTER TABLE experiment_run DROP COLUMN scheduled_end_time;
ALTER TABLE experiment_run ADD COLUMN cancellation_time TIMESTAMP WITH TIME ZONE;