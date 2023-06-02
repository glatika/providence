
-- +migrate Up
CREATE TABLE IF NOT EXISTS stock_variants (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    variant VARCHAR(512) NOT NULL,
	permissions TEXT NOT NULL,
	abilities TEXT NOT NULL,
	UNIQUE(variant)
);

-- +migrate Down
DROP TABLE IF EXISTS stock_variants;