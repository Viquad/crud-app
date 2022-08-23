CREATE TABLE IF NOT EXISTS users (
    id SERIAL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    registered_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(id, email)
);