CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL UNIQUE,
    user_id INT UNIQUE,
    balance INT DEFAULT 0 NOT NULL,
    currency VARCHAR(10) NOT NULL,
    last_update TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id),
    FOREIGN KEY (user_id)
        REFERENCES users(id)
);