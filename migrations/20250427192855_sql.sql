-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs (
    id           SERIAL        PRIMARY KEY,
    method_name  TEXT   NOT NULL,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    ctx        TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE logs;
-- +goose StatementEnd
