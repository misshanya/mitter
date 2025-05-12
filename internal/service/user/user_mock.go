package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/models"
)

var (
	testUser = models.User{
		ID:             testUserID,
		Login:          "testuser",
		Name:           "Test User",
		HashedPassword: "$2a$10$OW9yD0TyX0pOBO2MzJhtpeOC6O694OS37VJnnaJKFm.rUFt5fy4O6",
	}
	testUserID = uuid.MustParse("b096376a-5fa9-4130-907a-709c67008a65")

	testUser2 = models.User{
		ID:             testUser2ID,
		Login:          "testuser2",
		Name:           "Test User2",
		HashedPassword: "$2a$10$OW9yD0TyX0pOBO2MzJhtpeOC6O694OS37VJnnaJKFm.rUFt5fy4O6",
	}
	testUser2ID = uuid.MustParse("38386ffe-54ac-48be-9244-a5144b41a014")
)

// Mock User repo
type mockUserRepo struct{}

func (r *mockUserRepo) CreateUser(ctx context.Context, user *models.UserCreate) (uuid.UUID, error) {
	_ = ctx
	_ = user

	return testUserID, nil
}

func (r *mockUserRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	_ = ctx
	_ = login

	return &testUser, nil
}

func (r *mockUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	_ = ctx
	_ = id

	return &testUser, nil
}

func (r *mockUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_ = ctx
	_ = id

	return nil
}

func (r *mockUserRepo) UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdate) error {
	_ = ctx
	_ = id
	_ = user

	return nil
}

func (r *mockUserRepo) GetCurrentPasswordHash(ctx context.Context, id uuid.UUID) (string, error) {
	_ = ctx
	_ = id

	return testUser.HashedPassword, nil
}

func (r *mockUserRepo) ChangePassword(ctx context.Context, id uuid.UUID, newHashedPassword string) error {
	_ = ctx
	_ = id
	_ = newHashedPassword

	return nil
}

func (r *mockUserRepo) FollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error {
	_ = ctx
	_ = followerID
	_ = followeeID

	return nil
}

func (r *mockUserRepo) UnfollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) error {
	_ = ctx
	_ = followerID
	_ = followeeID

	return nil
}

func (r *mockUserRepo) GetUserFollows(ctx context.Context, followerID uuid.UUID, limit, offset int32) ([]uuid.UUID, error) {
	_ = ctx
	_ = followerID
	_ = limit
	_ = offset

	return []uuid.UUID{testUserID}, nil
}

func (r *mockUserRepo) GetUserFollowers(ctx context.Context, followeeID uuid.UUID, limit, offset int32) ([]uuid.UUID, error) {
	_ = ctx
	_ = followeeID
	_ = limit
	_ = offset

	return []uuid.UUID{testUserID}, nil
}

func (r *mockUserRepo) GetUserFriends(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]uuid.UUID, error) {
	_ = ctx
	_ = userID
	_ = limit
	_ = offset

	return []uuid.UUID{testUserID}, nil
}

// Mock user metrics
type mockUserMetrics struct {
	FakeUsersCount int
}

func (m *mockUserMetrics) AddUser() {
	m.FakeUsersCount++
}

func (m *mockUserMetrics) DeleteUser() {
	m.FakeUsersCount--
}
