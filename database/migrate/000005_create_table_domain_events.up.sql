BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS domain_event
(
    id      VARCHAR NOT NULL,
    version INTEGER NOT NULL,
    payload JSON    NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS id_version
    ON domain_event (id, version);

CREATE INDEX IF NOT EXISTS id
    ON domain_event (id);

ALTER TABLE domain_event
    ADD CONSTRAINT id_version
        UNIQUE USING INDEX id_version;

CREATE OR REPLACE RULE rule_domain_event_nodelete AS
    ON DELETE TO domain_event DO INSTEAD NOTHING;
CREATE OR REPLACE RULE rule_domain_event_noupdate AS
    ON UPDATE TO domain_event DO INSTEAD NOTHING;

CREATE TABLE IF NOT EXISTS domain_event_outbox
(
    id      SERIAL PRIMARY KEY,
    payload JSON NOT NULL
);

COMMIT TRANSACTION;
