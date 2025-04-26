package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type authRepo interface {
	GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error)
}

type AuthMiddleware struct {
	authRepo authRepo
}

func NewAuthMiddleware(authRepo authRepo) *AuthMiddleware {
	return &AuthMiddleware{authRepo: authRepo}
}

func (a *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization required")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header")
		}
		token := parts[1]

		userID, err := a.authRepo.GetUserIDByToken(c.Request().Context(), token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		c.Set("userID", userID)
		return next(c)
	}
}
