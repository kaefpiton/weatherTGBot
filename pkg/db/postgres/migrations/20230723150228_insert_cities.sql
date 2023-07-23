-- +goose Up
-- +goose StatementBegin
INSERT INTO cities (title,alias) VALUES('Moscow','Москва');
INSERT INTO cities (title,alias) VALUES('Rostov','Ростов');
INSERT INTO cities (title,alias) VALUES('Agalatovo','Алгатово');
INSERT INTO cities (title,alias) VALUES('Minsk','Минск');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE cities;
-- +goose StatementEnd
