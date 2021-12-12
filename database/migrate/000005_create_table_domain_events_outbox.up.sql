CREATE TABLE IF NOT EXISTS domain_event
(
    id         VARCHAR NOT NULL PRIMARY KEY,
    payload    jsonb   NOT NULL,
    event_type VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
