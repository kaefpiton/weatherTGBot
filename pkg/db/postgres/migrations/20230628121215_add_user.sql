-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    ID serial NOT NULL,
    firstname VARCHAR(256) NOT NULL,
    lastname VARCHAR(256) NOT NULL,
    chat_id int NOT NULL,
    last_usage timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(ID)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
