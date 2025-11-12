-- +goose Up
-- +goose StatementBegin
CREATE TABLE pull_requests_to_reviewers
(
    pull_request_id uuid NOT NULL,
    user_id         uuid NOT NULL,
    PRIMARY KEY (pull_request_id, user_id),
    FOREIGN KEY (pull_request_id) REFERENCES pull_requests (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pull_requests_to_reviewers;
-- +goose StatementEnd
