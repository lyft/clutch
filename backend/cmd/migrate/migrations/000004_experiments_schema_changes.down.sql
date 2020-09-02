DROP TABLE IF EXISTS experiment_config;
DROP TABLE IF EXISTS experiment_run;
CREATE TABLE IF NOT EXISTS experiments (
    id BIGINT,
    details JSONB
);
