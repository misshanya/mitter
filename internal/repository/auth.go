package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/models"
	"github.com/redis/go-redis/v9"
	"log/slog"
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
	userIDString, err := r.rdb.Get(ctx, token).Result()
	if err != nil {
		slog.Error("error getting user id from redis", slog.String("token", token))
		return uuid.Nil, err
	}

	// Parse string to uuid
	id, err := uuid.Parse(userIDString)
	if err != nil {
		slog.Error("error parsing user id from redis", slog.String("token", token), slog.String("userID", userIDString), slog.String("error", err.Error()))
		return uuid.Nil, err
	}
	return id, nil
}
