CREATE TABLE IF NOT EXISTS audit_events(
    id BIGSERIAL PRIMARY KEY,
    occurred_at TIMESTAMP WITH TIME ZONE,
    details JSONB
);
CREATE INDEX IF NOT EXISTS sort_audit_events ON audit_events (occurred_at);
CREATE INDEX IF NOT EXISTS audit_events_json ON audit_events USING GIN (details jsonb_path_ops)
