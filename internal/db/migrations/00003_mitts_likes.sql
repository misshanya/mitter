-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS mitts_likes (
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    mitt_id UUID NOT NULL REFERENCES mitts(id) ON DELETE CASCADE,
    liked_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mitts_likes;
-- +goose StatementEnd
