USE quotes;

INSERT INTO asset_categories (id, name) VALUES
    (1, 'stock'),
    (2, 'etf'),
    (3, 'crypto')
ON DUPLICATE KEY UPDATE name = VALUES(name);

INSERT INTO exchanges (id, code, name, currency, tz) VALUES
    (1, 'TWSE', 'Taiwan Stock Exchange', 'TWD', 'Asia/Taipei'),
    (2, 'BINANCE', 'Binance', 'USDT', 'UTC'),
    (3, 'Taiwan', 'Yahoo Finance Taiwan', 'TWD', 'Asia/Taipei'),
    (4, 'TWO', 'Taiwan OTC Exchange', 'TWD', 'Asia/Taipei'),
    (5, 'CCC', 'Yahoo Crypto Composite', 'USD', 'UTC')
ON DUPLICATE KEY UPDATE
    code = VALUES(code),
    name = VALUES(name),
    currency = VALUES(currency),
    tz = VALUES(tz);

