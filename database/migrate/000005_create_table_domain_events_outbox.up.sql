CREATE TABLE IF NOT EXISTS domain_events_outbox (
    id VARCHAR NOT NULL PRIMARY KEY,
    payload jsonb
)