-- +goose Up
-- +goose StatementBegin

CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description VARCHAR DEFAULT NULL,
    data JSONB NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE projects;

-- +goose StatementEnd