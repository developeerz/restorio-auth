CREATE TABLE IF NOT EXISTS auths (
    id VARCHAR(16) PRIMARY KEY,
    description TEXT NOT NULL
);

INSERT INTO auths (id, description) VALUES
('ADMIN', 'Права админа'),
('USER', 'Права пользователя');