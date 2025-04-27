package service

import (
	"context"
	"github.com/misshanya/mitter/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Tests
func TestUserService_GetUser(t *testing.T) {
	service := NewUserService(&mockUserRepo{})
	ctx := context.Background()

	user, err := service.GetUser(ctx, testUserID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testUser, *user)
}

func TestUserService_DeleteUser(t *testing.T) {
	service := NewUserService(&mockUserRepo{})
	ctx := context.Background()

	err := service.DeleteUser(ctx, testUserID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	service := NewUserService(&mockUserRepo{})
	ctx := context.Background()

	newName := "new name"
	userUpdate := &models.UserUpdate{
		Name: &newName,
	}
	err := service.UpdateUser(ctx, testUserID, userUpdate)
	if err != nil {
		t.Fatal(err)
	}
}
