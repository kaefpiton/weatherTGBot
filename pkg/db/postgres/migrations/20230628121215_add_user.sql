-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    users_id serial NOT NULL,
    users_firstname VARCHAR(40) NOT NULL,
    users_lastname VARCHAR(40) NOT NULL,
    users_chatid int NOT NULL,
    users_date_of_last_usage timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(users_id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
