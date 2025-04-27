package service

import (
	"context"
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
