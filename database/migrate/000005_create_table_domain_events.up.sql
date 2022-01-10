BEGIN TRANSACTION READ WRITE;

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

COMMIT TRANSACTION;
