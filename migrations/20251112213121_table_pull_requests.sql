-- +goose Up
-- +goose StatementBegin
CREATE TABLE pull_requests
(
    id                  uuid PRIMARY KEY    NOT NULL,
    name                text                NOT NULL,
    author_id           uuid                NOT NULL,
    status              pull_request_status NOT NULL,
    need_more_reviewers bool                NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pull_requests;
-- +goose StatementEnd
