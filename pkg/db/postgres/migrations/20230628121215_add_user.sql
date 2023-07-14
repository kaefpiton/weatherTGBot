-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    chat_id int NOT NULL,
    firstname VARCHAR(256) NOT NULL,
    lastname VARCHAR(256) NOT NULL,
    state VARCHAR(256) NOT NULL DEFAULT 'auth',
    last_usage timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(chat_id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
