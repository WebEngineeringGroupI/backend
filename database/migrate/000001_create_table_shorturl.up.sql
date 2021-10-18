CREATE TABLE IF NOT EXISTS shorturl
(
    hash       VARCHAR(8)     NOT NULL PRIMARY KEY,
    long_url   VARCHAR(65535) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP
);
