package auth

import (
	"context"
	"testing"

	"github.com/misshanya/mitter/internal/models"
	"github.com/stretchr/testify/assert"
)

// Tests
func TestAuthService_SignIn(t *testing.T) {
	service := NewAuthService(&mockUserRepo{}, &mockAuthRepo{}, &mockUserMetrics{})

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
	service := NewAuthService(&mockUserRepo{}, &mockAuthRepo{}, &mockUserMetrics{})

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

func TestAuthService_ChangePassword(t *testing.T) {
	service := NewAuthService(&mockUserRepo{}, &mockAuthRepo{}, &mockUserMetrics{})

	ctx := context.Background()

	changPwd := &models.ChangePassword{
		OldPassword: "qwerty123456",
		NewPassword: "qwerty1234567",
	}
	err := service.ChangePassword(ctx, testUserID, changPwd)
	if err != nil {
		t.Fatal(err)
	}
}
