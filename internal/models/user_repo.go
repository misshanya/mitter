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

	FollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error
	UnfollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error
	GetUserFollows(ctx context.Context, followerID uuid.UUID, limit, offset int32) ([]uuid.UUID, error)
	GetUserFollowers(ctx context.Context, followeeID uuid.UUID, limit, offset int32) ([]uuid.UUID, error)
	GetUserFriends(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}
