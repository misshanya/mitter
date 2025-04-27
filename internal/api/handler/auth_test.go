package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misshanya/mitter/internal/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock service
type mockAuthService struct{}

func (s *mockAuthService) SignIn(ctx context.Context, creds models.SignIn) (string, *models.HTTPError) {
	_ = ctx
	_ = creds

	return "8a67006c-692f-4e75-b547-84a46707a5cb", nil
}

func (s *mockAuthService) SignUp(ctx context.Context, user *models.UserCreate) (uuid.UUID, *models.HTTPError) {
	_ = ctx
	_ = user

	return uuid.MustParse("b096376a-5fa9-4130-907a-709c67008a65"), nil
}

func (s *mockAuthService) ChangePassword(ctx context.Context, id uuid.UUID, changePassword *models.ChangePassword) *models.HTTPError {
	_ = ctx
	_ = id
	_ = changePassword

	return nil
}

// Mock auth middleware
func mockRequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("userID", "b096376a-5fa9-4130-907a-709c67008a65")
		return next(c)
	}
}

// Tests
func TestAuthHandler_SignIn(t *testing.T) {
	e := echo.New()
	handler := NewAuthHandler(&mockAuthService{}, mockRequireAuth)

	g := e.Group("/api/v1/auth")
	handler.Routes(g)

	// Create request
	reqBody := `{"login":"testuser","password":"qwerty123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	if assert.NoError(t, handler.signIn(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResp := `{"token": "8a67006c-692f-4e75-b547-84a46707a5cb"}`
		assert.JSONEq(t, expectedResp, rec.Body.String())
	}
}

func TestAuthHandler_SignUp(t *testing.T) {
	e := echo.New()
	handler := NewAuthHandler(&mockAuthService{}, mockRequireAuth)

	g := e.Group("/api/v1/auth")
	handler.Routes(g)

	// Create request
	reqBody := `{"login":"testuser","name":"Test User","password":"qwerty123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	if assert.NoError(t, handler.signUp(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		exceptedResp := `{"id": "b096376a-5fa9-4130-907a-709c67008a65"}`
		assert.JSONEq(t, exceptedResp, rec.Body.String())
	}
}
