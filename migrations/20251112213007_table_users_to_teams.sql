-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_to_teams
(
    user_id uuid NOT NULL,
    team_id uuid NOT NULL,
    PRIMARY KEY (user_id, team_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_to_teams;
-- +goose StatementEnd
