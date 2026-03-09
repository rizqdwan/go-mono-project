CREATE TABLE IF NOT EXISTS departments (
    id         BIGSERIAL   PRIMARY KEY,
    label      VARCHAR(10) NOT NULL,
    name       VARCHAR(60) NOT NULL,
    group_id   BIGINT      REFERENCES groups(id) ON DELETE SET NULL,
    created_at TIMESTAMP   DEFAULT NOW(),
    updated_at TIMESTAMP   DEFAULT NOW()
);

CREATE INDEX idx_departments_group_id ON departments(group_id);