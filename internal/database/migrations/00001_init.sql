-- +goose Up
-- +goose StatementBegin
CREATE TABLE info(
    id SERIAL PRIMARY KEY,
    value TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE info;
-- +goose StatementEnd
