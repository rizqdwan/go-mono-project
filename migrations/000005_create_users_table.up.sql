CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL    PRIMARY KEY,
    email         VARCHAR(50)  NOT NULL UNIQUE,
    name          VARCHAR(50)  NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role_id       BIGINT       REFERENCES user_roles(id) ON DELETE SET NULL,
    department_id BIGINT       REFERENCES departments(id) ON DELETE SET NULL,
    position_id   VARCHAR(10)  REFERENCES user_positions(id) ON DELETE SET NULL,
    created_at    TIMESTAMP    DEFAULT NOW(),
    updated_at    TIMESTAMP    DEFAULT NOW()
);

CREATE INDEX idx_users_role_id       ON users(role_id);
CREATE INDEX idx_users_department_id ON users(department_id);
CREATE INDEX idx_users_position_id   ON users(position_id);