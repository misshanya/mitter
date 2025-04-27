package models

import (
	"context"
	"github.com/google/uuid"
)

type AuthRepository interface {
	SaveToken(ctx context.Context, token *Token) error
	GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error)
}
