-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stickers (
    ID serial NOT NULL,
    title VARCHAR(256) NOT NULL,
    code VARCHAR(256) NOT NULL,
    sticker_type VARCHAR(256) NOT NULL,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(ID)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table stickers;
-- +goose StatementEnd
