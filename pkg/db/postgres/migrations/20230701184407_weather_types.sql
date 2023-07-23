-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS weather_types (
    ID serial NOT NULL,
    title VARCHAR(256) NOT NULL,
    alias VARCHAR(256) NOT NULL,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(ID)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table weather_types;
-- +goose StatementEnd
