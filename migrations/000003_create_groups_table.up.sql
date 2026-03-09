CREATE TABLE IF NOT EXISTS groups (
    id         BIGSERIAL   PRIMARY KEY,
    label      VARCHAR(10),
    name       VARCHAR(60),
    created_at TIMESTAMP   DEFAULT NOW(),
    updated_at TIMESTAMP   DEFAULT NOW()
);