-- +goose Up
-- +goose StatementBegin
-- status ENUM
CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE IF NOT EXISTS orders
(
    id         SERIAL PRIMARY KEY,
    number     VARCHAR(255) UNIQUE NOT NULL,
    accrual    DECIMAL(10, 2),
    user_id    INT                 NOT NULL REFERENCES users (id),
    status     order_status        NOT NULL DEFAULT 'NEW',
    uploaded_at TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
-- +goose StatementEnd
