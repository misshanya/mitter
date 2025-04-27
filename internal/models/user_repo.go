package models

import (
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *UserCreate) (uuid.UUID, error)
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)

	DeleteUser(ctx context.Context, id uuid.UUID) error

	UpdateUser(ctx context.Context, id uuid.UUID, user *UserUpdate) error

	GetCurrentPasswordHash(ctx context.Context, id uuid.UUID) (string, error)
	ChangePassword(ctx context.Context, id uuid.UUID, newHashedPassword string) error
}
