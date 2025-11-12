-- +goose Up
-- +goose StatementBegin
CREATE TYPE pull_request_status AS ENUM ('OPEN', 'MERGED');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE pull_request_status;
-- +goose StatementEnd
