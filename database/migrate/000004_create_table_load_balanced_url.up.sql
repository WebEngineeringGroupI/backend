CREATE TABLE IF NOT EXISTS load_balanced_url
(
    hash         VARCHAR(8)     NOT NULL,
    original_url VARCHAR(65535) NOT NULL,
    is_valid     BOOLEAN   DEFAULT FALSE,
    created_at   TIMESTAMP DEFAULT now(),
    updated_at   TIMESTAMP,
    CONSTRAINT hash_original_url PRIMARY KEY (hash, original_url)
);
