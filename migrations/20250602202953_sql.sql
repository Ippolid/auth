-- +goose Up
-- +goose StatementBegin
CREATE TABLE access (
    id SERIAL PRIMARY KEY,
    role BOOL NOT NULL,
    endpoint TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE access;
-- +goose StatementEnd
