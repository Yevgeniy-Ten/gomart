-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS balance
(
    id        SERIAL PRIMARY KEY,
    user_id   INT NOT NULL REFERENCES users (id),
    current   INT NOT NULL DEFAULT 0 CHECK ( current >= 0 ),
    withdrawn INT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS balance;
-- +goose StatementEnd
