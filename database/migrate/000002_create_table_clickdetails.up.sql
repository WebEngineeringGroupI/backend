CREATE TABLE IF NOT EXISTS clickdetails
(
    id         SERIAL PRIMARY KEY,
    hash       VARCHAR(8)  NOT NULL,
    ip         VARCHAR(39) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP
);
