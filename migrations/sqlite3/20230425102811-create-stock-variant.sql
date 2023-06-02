
-- +migrate Up
CREATE TABLE IF NOT EXISTS stock_variants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    variant TEXT NOT NULL UNIQUE,
    permissions TEXT NOT NULL,
    abilities TEXT NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS stock_variants;