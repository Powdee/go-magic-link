CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    magic_token VARCHAR(255),
    token_expiration TIMESTAMP
);

