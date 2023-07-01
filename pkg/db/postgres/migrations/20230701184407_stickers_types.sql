-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sticker_types (
    ID serial NOT NULL,
    title VARCHAR(256) NOT NULL,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(ID)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table sticker_types;
-- +goose StatementEnd
