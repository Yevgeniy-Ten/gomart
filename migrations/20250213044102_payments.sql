-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payments
(
    id           SERIAL PRIMARY KEY,
    user_id      INT            NOT NULL REFERENCES users (id),
    "order"      VARCHAR(255)   NOT NULL,
    sum          DECIMAL(10, 2) NOT NULL,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
-- +goose StatementEnd
