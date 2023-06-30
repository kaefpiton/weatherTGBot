-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stickers (
    stickers_id serial NOT NULL,
    stickers_name VARCHAR(256) NOT NULL,
    stickers_code VARCHAR(256) NOT NULL,
    stickers_type_id VARCHAR(256) NOT NULL,
    PRIMARY KEY(stickers_id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table stickers;
-- +goose StatementEnd
