package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/misshanya/mitter/internal/db/sqlc/storage"
	"github.com/misshanya/mitter/internal/models"
)

type UserRepository struct {
	queries *storage.Queries
}

func NewUserRepository(q *storage.Queries) *UserRepository {
	return &UserRepository{queries: q}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.UserCreate) (uuid.UUID, error) {
	return r.queries.CreateUser(ctx, storage.CreateUserParams{
		Login:          user.Login,
		Name:           user.Name,
		Hashedpassword: user.HashedPassword,
	})
}

func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	userDB, err := r.queries.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	// Convert db user to my model
	user := &models.User{
		ID:             userDB.ID,
		Login:          userDB.Login,
		Name:           userDB.Name,
		HashedPassword: userDB.Password,
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	userDB, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert db user to my model
	user := &models.User{
		ID:             userDB.ID,
		Login:          userDB.Login,
		Name:           userDB.Name,
		HashedPassword: userDB.Password,
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdate) error {
	name := pgtype.Text{}
	if user.Name != nil {
		name = pgtype.Text{String: *user.Name, Valid: true}
	}

	return r.queries.UpdateUser(ctx, storage.UpdateUserParams{
		Name: name,
		ID:   id,
	})
}

func (r *UserRepository) GetCurrentPasswordHash(ctx context.Context, id uuid.UUID) (string, error) {
	return r.queries.GetCurrentPasswordHash(ctx, id)
}

func (r *UserRepository) ChangePassword(ctx context.Context, id uuid.UUID, newHashedPassword string) error {
	return r.queries.UpdatePassword(ctx, storage.UpdatePasswordParams{
		Password: newHashedPassword,
		ID:       id,
	})
}

func (r *UserRepository) FollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error {
	return r.queries.FollowUser(ctx, storage.FollowUserParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
}

func (r *UserRepository) UnfollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error {
	return r.queries.UnfollowUser(ctx, storage.UnfollowUserParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
}

func (r *UserRepository) GetUserFollows(ctx context.Context, followerID uuid.UUID) ([]uuid.UUID, error) {
	return r.queries.GetUserFollows(ctx, followerID)
}

func (r *UserRepository) GetUserFollowers(ctx context.Context, followeeID uuid.UUID) ([]uuid.UUID, error) {
	return r.queries.GetUserFollowers(ctx, followeeID)
}
