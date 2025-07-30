-- +goose Up
-- +goose StatementBegin
CREATE TABLE savings(
    UUID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    target_amount DECIMAL(10, 2) NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    image VARCHAR(255) NOT NULL,
    filling_plan VARCHAR(7) NOT NULL CHECK (filling_plan IN ('Daily', 'Weekly', 'Monthly')),
    filling_nominal DECIMAL(10, 2) NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (user_uuid) REFERENCES users(uuid)
);

-- Create indexing
CREATE INDEX idx_savings_user_uuid ON savings(user_uuid);
CREATE INDEX idx_savings_currency_code ON savings(currency_code);
CREATE INDEX idx_savings_filling_plan ON savings(filling_plan);
CREATE INDEX idx_savings_deleted_at ON savings(deleted_at);
CREATE INDEX idx_savings_name ON savings(name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS savings;
DROP INDEX IF EXISTS idx_savings_user_uuid;
DROP INDEX IF EXISTS idx_savings_currency_code;
DROP INDEX IF EXISTS idx_savings_filling_plan;
DROP INDEX IF EXISTS idx_savings_deleted_at;
DROP INDEX IF EXISTS idx_savings_name;
-- +goose StatementEnd
