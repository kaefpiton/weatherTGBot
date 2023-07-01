-- +goose Up
-- +goose StatementBegin
INSERT INTO sticker_types (title) VALUES('high temperature');
INSERT INTO sticker_types (title) VALUES('normal temperature');
INSERT INTO sticker_types (title) VALUES('cold temperature');
INSERT INTO sticker_types (title) VALUES('frost temperature');
INSERT INTO sticker_types (title) VALUES('pressure high');
INSERT INTO sticker_types (title) VALUES('pressure normal');
INSERT INTO sticker_types (title) VALUES('high wind');
INSERT INTO sticker_types (title) VALUES('normal wind');
INSERT INTO sticker_types (title) VALUES('low wind');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE sticker_types;
-- +goose StatementEnd
