-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255),
    photo VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE UNIQUE INDEX idx_users_email ON users(email);

CREATE UNIQUE INDEX idx_users_phone_number ON users(phone_number);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);

CREATE INDEX idx_users_name ON users(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_users_phone_number;

DROP INDEX IF EXISTS idx_users_deleted_at;

DROP INDEX IF EXISTS idx_users_name;
-- +goose StatementEnd
