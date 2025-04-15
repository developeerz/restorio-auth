CREATE TABLE IF NOT EXISTS user_auths (
    user_telegram_id INTEGER REFERENCES users (telegram_id) ON DELETE CASCADE,
    auth_id VARCHAR(15) REFERENCES auths (id),
    PRIMARY KEY(user_telegram_id, auth_id)
);