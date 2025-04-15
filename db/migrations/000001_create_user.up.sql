CREATE TABLE IF NOT EXISTS users (
    telegram_id BIGINT PRIMARY KEY,
    firstname VARCHAR(63) NOT NULL,
    lastname VARCHAR(63) NOT NULL,
    telegram VARCHAR(63) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    registration_date TIMESTAMP NOT NULL
);