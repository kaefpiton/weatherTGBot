-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cities (
    ID serial NOT NULL,
    title VARCHAR(256) NOT NULL,
    alias VARCHAR(256) NOT NULL,
    clicks INT NOT NULL DEFAULT 0,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(ID)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table cities;
-- +goose StatementEnd