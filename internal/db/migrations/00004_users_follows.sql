-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_follows (
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    follower_id UUID NOT NULL REFERENCES users(id),
    followee_id UUID NOT NULL REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_users_follows_follower_id ON users_follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_users_follows_followee_id ON users_follows(followee_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_follows_follower_id;
DROP INDEX IF EXISTS idx_users_follows_followee_id;

DROP TABLE IF EXISTS users_follows;
-- +goose StatementEnd
