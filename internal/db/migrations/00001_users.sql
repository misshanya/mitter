-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id UUID NOT NULL default gen_random_uuid() PRIMARY KEY,
  login VARCHAR(50) NOT NULL UNIQUE,
  name TEXT NOT NULL,
  password TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
