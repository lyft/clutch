ALTER TABLE audit_events
    ADD COLUMN IF NOT EXISTS "sent" BOOLEAN DEFAULT FALSE;
CREATE INDEX IF NOT EXISTS sent_audit_events ON audit_events (sent) WHERE sent = FALSE;
