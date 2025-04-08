CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(63) NOT NULL,
    lastname VARCHAR(63) NOT NULL,
    telegram VARCHAR(63) UNIQUE NOT NULL,
    telegram_id BIGINT DEFAULT NULL,
    password VARCHAR(255) NOT NULL,
    registration_date TIMESTAMP NOT NULL
);