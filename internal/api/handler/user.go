package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misshanya/mitter/internal/models"
	"net/http"
)

type userService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type UserHandler struct {
	service  userService
	validate *validator.Validate
}

func NewUserHandler(service userService) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *UserHandler) Routes(group *echo.Group) {
	group.GET("/", h.GetMe)
}

func (h *UserHandler) GetMe(c echo.Context) error {
	return c.String(http.StatusOK, "TODO")
}
