package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/db/sqlc/storage"
	"github.com/misshanya/mitter/internal/models"
)

type UserRepository struct {
	queries *storage.Queries
}

func NewUserRepository(q *storage.Queries) *UserRepository {
	return &UserRepository{queries: q}
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.UserCreate) (uuid.UUID, error) {
	return r.queries.CreateUser(ctx, storage.CreateUserParams{
		Login: user.Login,
		Name:  user.Name,
	})
}
