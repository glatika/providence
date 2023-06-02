
-- +migrate Up
CREATE TABLE IF NOT EXISTS stocks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hwid TEXT NOT NULL,
    operating_system TEXT NOT NULL,
    variant INTEGER NOT NULL,
    'signature' TEXT NULL,
    FOREIGN KEY(variant) REFERENCES stock_variants(id)
);

-- +migrate Down
DROP TABLE IF EXISTS stocks;