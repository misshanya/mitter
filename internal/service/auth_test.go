package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
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

// Mock repos
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

// Tests
func TestAuthService_SignIn(t *testing.T) {
	service := NewAuthService(&mockUserRepo{}, &mockAuthRepo{})

	ctx := context.Background()

	creds := models.SignIn{
		Login:    testUser.Login,
		Password: "qwerty123456",
	}
	token, err := service.SignIn(ctx, creds)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.NotEmpty(t, token) {
		t.Fatal()
	}
}

func TestAuthService_SignUp(t *testing.T) {
	service := NewAuthService(&mockUserRepo{}, &mockAuthRepo{})

	ctx := context.Background()

	user := &models.UserCreate{
		Login:    testUser.Login,
		Name:     testUser.Name,
		Password: "qwerty123456",
	}
	id, err := service.SignUp(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.NotEmpty(t, id) {
		t.Fatal()
	}
}
