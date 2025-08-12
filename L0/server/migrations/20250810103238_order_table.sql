-- +goose Up
-- +goose StatementBegin
CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE payments (
    transaction VARCHAR(255) PRIMARY KEY,
    request_id VARCHAR(255),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(50) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL
);

CREATE TABLE items (
    chrt_id INTEGER PRIMARY KEY,
    track_number VARCHAR(50) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(10) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL
);

CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(255) NOT NULL,
    delivery_id INTEGER REFERENCES deliveries(id) NOT NULL,
    payment_transaction VARCHAR(255) REFERENCES payments(transaction) NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255) NOT NULL,
    shardkey VARCHAR(255) NOT NULL,
    sm_id INTEGER NOT NULL,
    date_created TIMESTAMP WITH TIME ZONE NOT NULL,
    oof_shard VARCHAR(255) NOT NULL
);

CREATE TABLE order_items (
    order_uid VARCHAR(255) REFERENCES orders(order_uid) NOT NULL,
    item_chrt_id INTEGER REFERENCES items(chrt_id) NOT NULL,
    PRIMARY KEY (order_uid, item_chrt_id)
);

CREATE INDEX idx_orders_order_uid ON orders(order_uid);
CREATE INDEX idx_deliveries_id ON deliveries(id);
CREATE INDEX idx_payments_transaction ON payments(transaction);
CREATE INDEX idx_items_chrt_id ON items(chrt_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_items_chrt_id;
DROP INDEX IF EXISTS idx_payments_transaction;
DROP INDEX IF EXISTS idx_deliveries_id;
DROP INDEX IF EXISTS idx_orders_order_uid;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;
-- +goose StatementEnd