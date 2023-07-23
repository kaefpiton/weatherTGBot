-- +goose Up
-- +goose StatementBegin
INSERT INTO weather_types (title, alias) VALUES('high temperature', 'жара');
INSERT INTO weather_types (title, alias) VALUES('normal temperature', 'нормальная температура');
INSERT INTO weather_types (title, alias) VALUES('cold temperature', 'холод');
INSERT INTO weather_types (title, alias) VALUES('frost temperature', 'дубак');
INSERT INTO weather_types (title, alias) VALUES('pressure high', 'высокое давление');
INSERT INTO weather_types (title, alias) VALUES('pressure normal', 'нормальное давление');
INSERT INTO weather_types (title, alias) VALUES('high wind', 'сильный втер');
INSERT INTO weather_types (title, alias) VALUES('normal wind', 'нормальный ветер');
INSERT INTO weather_types (title, alias) VALUES('low wind', 'слабый ветер');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE weather_types;
-- +goose StatementEnd
