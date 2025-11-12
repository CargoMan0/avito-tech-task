-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id        uuid PRIMARY KEY NOT NULL,
    name      text             NOT NULL,
    is_active bool             NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
