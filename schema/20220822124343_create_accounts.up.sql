CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    balance INT DEFAULT 0 NOT NULL,
    currency VARCHAR(10) NOT NULL,
    last_update TIMESTAMP DEFAULT NOW()        
);