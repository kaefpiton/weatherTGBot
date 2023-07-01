-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stickers (
    ID serial NOT NULL,
    title VARCHAR(256) NOT NULL,
    code VARCHAR(256) NOT NULL,
    sticker_type VARCHAR(256) NOT NULL,
    type_id INT NOT NULL,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(ID),
    CONSTRAINT fk_stickers_types FOREIGN KEY(type_id) REFERENCES sticker_types(ID)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table stickers;
-- +goose StatementEnd