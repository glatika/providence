
-- +migrate Up
CREATE TABLE IF NOT EXISTS stock_details (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    stock_id INTEGER NOT NULL UNIQUE,
    granted_permission TEXT NOT NULL,
    hostname TEXT NULL,
    country TEXT NULL,
    misc_detail TEXT NULL,
   FOREIGN KEY(stock_id) REFERENCES stocks(id)
       ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS stock_details;