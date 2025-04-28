-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS mitts (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    author UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mitts;
-- +goose StatementEnd
