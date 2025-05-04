package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misshanya/mitter/internal/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock service
type mockUserService struct{}

func (s *mockUserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, *models.HTTPError) {
	_ = ctx
	return &models.User{
		ID:             id,
		Login:          "testuser",
		Name:           "Test User",
		HashedPassword: "abracadabra",
	}, nil
}

func (s *mockUserService) DeleteUser(ctx context.Context, id uuid.UUID) *models.HTTPError {
	_ = ctx
	_ = id
	return nil
}

func (s *mockUserService) UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdate) *models.HTTPError {
	_ = ctx
	_ = id
	_ = user

	return nil
}

// Tests
func TestUserHandler_GetMe(t *testing.T) {
	e := echo.New()

	handler := NewUserHandler(&mockUserService{})

	g := e.Group("/api/v1/user")
	handler.Routes(g)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.Set("userID", uuid.MustParse("b096376a-5fa9-4130-907a-709c67008a65"))

	if assert.NoError(t, handler.getMe(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "testuser")
		assert.Contains(t, rec.Body.String(), "Test User")
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	e := echo.New()

	handler := NewUserHandler(&mockUserService{})

	g := e.Group("/api/v1/user")
	handler.Routes(g)

	// Create request
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/user/", nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.Set("userID", uuid.MustParse("b096376a-5fa9-4130-907a-709c67008a65"))

	if assert.NoError(t, handler.deleteUser(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
