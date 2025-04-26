package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type AuthRepository struct {
	rdb *redis.Client
}

func NewAuthRepository(rdb *redis.Client) *AuthRepository {
	return &AuthRepository{
		rdb: rdb,
	}
}

func (r *AuthRepository) SaveToken(ctx context.Context, token *models.Token) error {
	return r.rdb.Set(ctx, token.Token.String(), token.UserID.String(), 24*time.Hour).Err()
}

func (r *AuthRepository) GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error) {
	// Get uuid string
	cmd := r.rdb.Get(ctx, token)
	if cmd.Err() != nil {
		return uuid.Nil, cmd.Err()
	}

	// Parse string to uuid
	id, err := uuid.Parse(cmd.String())
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
