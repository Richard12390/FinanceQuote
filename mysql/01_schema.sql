USE quotes;

CREATE TABLE IF NOT EXISTS asset_categories (
    id TINYINT PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS exchanges (
    id SMALLINT PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    currency VARCHAR(16) NOT NULL,
    tz VARCHAR(128) NULL
);


CREATE TABLE IF NOT EXISTS symbols (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    source VARCHAR(32) NOT NULL,
    category_id TINYINT NOT NULL,
    exchange_id SMALLINT NOT NULL,
    symbol VARCHAR(64) NOT NULL,
    symbol_origin VARCHAR(64) NOT NULL,
    currency VARCHAR(16) NOT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    UNIQUE KEY uq_src_cat_sym (source, category_id, symbol),
    CONSTRAINT fk_symbol_category FOREIGN KEY (category_id) REFERENCES asset_categories (id),
    CONSTRAINT fk_symbol_exchange FOREIGN KEY (exchange_id) REFERENCES exchanges (id)
);

CREATE TABLE favorites (
    user_id BIGINT NOT NULL,
    symbol  VARCHAR(64) NOT NULL,
    PRIMARY KEY (user_id, symbol)
);

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account       VARCHAR(64)  NOT NULL UNIQUE,
    password      VARCHAR(255) NOT NULL,
    display_name  VARCHAR(128) NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ALTER TABLE symbols DROP FOREIGN KEY fk_symbol_exchange;
-- ALTER TABLE exchanges
--     MODIFY COLUMN id SMALLINT NOT NULL,
--     MODIFY COLUMN code VARCHAR(32) NOT NULL UNIQUE,
--     MODIFY COLUMN name VARCHAR(128) NOT NULL,
--     MODIFY COLUMN currency VARCHAR(16) NOT NULL,
--     MODIFY COLUMN tz VARCHAR(128) NULL;
-- ALTER TABLE symbols
--     MODIFY COLUMN currency VARCHAR(16) NOT NULL,
--     ADD CONSTRAINT fk_symbol_exchange FOREIGN KEY (exchange_id) REFERENCES exchanges (id);

-- CREATE TABLE IF NOT EXISTS daily_bars (
--     symbol_id BIGINT NOT NULL,
--     trade_date DATE NOT NULL,
--     open DECIMAL(18, 8) NOT NULL,
--     high DECIMAL(18, 8) NOT NULL,
--     low DECIMAL(18, 8) NOT NULL,
--     close DECIMAL(18, 8) NOT NULL,
--     adj_close DECIMAL(18, 8) NULL,
--     volume DECIMAL(18, 8) NULL,
--     vwap DECIMAL(18, 8) NULL,
--     raw JSON NULL,
--     PRIMARY KEY (symbol_id, trade_date),
--     CONSTRAINT fk_quote_symbol FOREIGN KEY (symbol_id) REFERENCES symbols (id)
-- );
