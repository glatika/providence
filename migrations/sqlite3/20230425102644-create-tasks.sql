-- +migrate Up
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stock_id INTEGER NOT NULL,
    instruction TEXT NOT NULL,
    argument TEXT NULL,
    delivered BOOLEAN NOT NULL DEFAULT False,
    delivered_at DATETIME  NULL,
    reported BOOLEAN NOT NULL DEFAULT False,
    reported_at DATETIME  NULL,
    success BOOLEAN NULL,
    report TEXT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS tasks;