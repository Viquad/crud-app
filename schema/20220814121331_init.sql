-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    balance INT DEFAULT 0 NOT NULL,
    currency VARCHAR(10) NOT NULL,
    last_update TIMESTAMP DEFAULT NOW()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts;

-- +goose StatementEnd