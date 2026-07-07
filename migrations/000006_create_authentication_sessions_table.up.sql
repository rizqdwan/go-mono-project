CREATE TABLE IF NOT EXISTS authentication_sessions (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT,
    is_active     BOOLEAN   DEFAULT TRUE,
    browser_info  TEXT,
    created_at    TIMESTAMP DEFAULT NOW(),
    last_activity TIMESTAMP DEFAULT NOW(),
    expires_at    TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_auth_sessions_user_id            ON authentication_sessions(user_id);
CREATE INDEX idx_auth_sessions_user_id_is_active  ON authentication_sessions(user_id, is_active);