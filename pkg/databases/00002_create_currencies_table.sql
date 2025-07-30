-- +goose Up
-- +goose StatementBegin
CREATE TABLE currencies(
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_name VARCHAR(255) NOT NULL,
    currency_name VARCHAR(255) NOT NULL,
    country_flag VARCHAR(10) NOT NULL, -- Emoji flag (🇮🇩, 🇺🇸, etc)
    currency_symbol VARCHAR(10) NOT NULL, -- Symbol ($, €, £, etc)
    currency_code VARCHAR(3) NOT NULL, -- ISO code (USD, EUR, IDR, etc)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- Create indexes for better performance
CREATE UNIQUE INDEX idx_currencies_currency_code ON currencies(currency_code);
CREATE INDEX idx_currencies_country_name ON currencies(country_name);
CREATE INDEX idx_currencies_currency_name ON currencies(currency_name);
CREATE INDEX idx_currencies_deleted_at ON currencies(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS currencies;

DROP INDEX IF EXISTS idx_currencies_currency_code;
DROP INDEX IF EXISTS idx_currencies_country_name;
DROP INDEX IF EXISTS idx_currencies_currency_name;
DROP INDEX IF EXISTS idx_currencies_deleted_at;
-- +goose StatementEnd
