package user

import (
	"context"
	"github.com/misshanya/mitter/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Tests
func TestUserService_GetUser(t *testing.T) {
	service := NewUserService(&mockUserRepo{}, &mockUserMetrics{})
	ctx := context.Background()

	user, err := service.GetUser(ctx, testUserID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testUser, *user)
}

func TestUserService_DeleteUser(t *testing.T) {
	service := NewUserService(&mockUserRepo{}, &mockUserMetrics{})
	ctx := context.Background()

	err := service.DeleteUser(ctx, testUserID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	service := NewUserService(&mockUserRepo{}, &mockUserMetrics{})
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

func TestUserService_FollowUser(t *testing.T) {
	service := NewUserService(&mockUserRepo{}, &mockUserMetrics{})
	ctx := context.Background()

	err := service.FollowUser(ctx, testUserID, testUser2ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserService_UnfollowUser(t *testing.T) {
	service := NewUserService(&mockUserRepo{}, &mockUserMetrics{})
	ctx := context.Background()

	err := service.UnfollowUser(ctx, testUserID, testUser2ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserService_GetUserFollows(t *testing.T) {
	service := NewUserService(&mockUserRepo{}, &mockUserMetrics{})
	ctx := context.Background()

	follows, err := service.GetUserFollows(ctx, testUserID, 30, 0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []*models.User{&testUser}, follows)
}

func TestUserService_GetUserFollowers(t *testing.T) {
	service := NewUserService(&mockUserRepo{}, &mockUserMetrics{})
	ctx := context.Background()

	followers, err := service.GetUserFollowers(ctx, testUserID, 30, 0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []*models.User{&testUser}, followers)
}
