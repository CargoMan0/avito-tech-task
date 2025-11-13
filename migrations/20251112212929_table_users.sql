-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id        uuid PRIMARY KEY NOT NULL,
    name      text             NOT NULL,
    is_active bool             NOT NULL,
    team_id   uuid             NOT NULL,
    FOREIGN KEY (team_id) REFERENCES teams (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd