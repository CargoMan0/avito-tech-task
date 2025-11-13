-- +goose Up
-- +goose StatementBegin
CREATE TABLE teams
(
    id   uuid PRIMARY KEY NOT NULL,
    name text UNIQUE      NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE teams;
-- +goose StatementEnd
