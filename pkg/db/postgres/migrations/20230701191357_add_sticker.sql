-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stickers (
    code VARCHAR(256) NOT NULL UNIQUE,
    title VARCHAR(256) NOT NULL,
    type_id INT NOT NULL,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(code),
    CONSTRAINT fk_stickers_types FOREIGN KEY(type_id) REFERENCES weather_types(ID)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table stickers;
-- +goose StatementEnd