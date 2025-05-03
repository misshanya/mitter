package auth

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
)

// Mock Auth repo
type mockAuthRepo struct{}

func (r *mockAuthRepo) SaveToken(ctx context.Context, token *models.Token) error {
	_ = ctx
	_ = token

	return nil
}

func (r *mockAuthRepo) GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error) {
	_ = ctx
	_ = token

	return testUserID, nil
}

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
