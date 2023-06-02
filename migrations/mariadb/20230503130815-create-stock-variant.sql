
-- +migrate Up
CREATE TABLE IF NOT EXISTS stocks (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    hwid TEXT NOT NULL,
    operating_system TEXT NOT NULL,
    variant INTEGER NOT NULL,
    signature TEXT NULL,
    FOREIGN KEY(variant)
     	REFERENCES stock_variant(id)
);

-- INFO: BLAME Planetscale
-- ALTER TABLE stocks
-- ADD CONSTRAINT stock_variant_rl
-- 

-- +migrate Down
DROP TABLE IF EXISTS stocks;