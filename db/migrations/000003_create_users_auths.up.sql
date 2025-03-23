CREATE TABLE IF NOT EXISTS user_auths (
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    auth_id VARCHAR(15) REFERENCES auths (id),
    PRIMARY KEY(user_id, auth_id)
);