package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misshanya/mitter/internal/api/dto"
	"github.com/misshanya/mitter/internal/models"
	"net/http"
)

type userService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, *models.HTTPError)
	DeleteUser(ctx context.Context, id uuid.UUID) *models.HTTPError

	UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdate) *models.HTTPError
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
	group.DELETE("/", h.DeleteUser)
	group.PATCH("/", h.UpdateUser)
}

func (h *UserHandler) GetMe(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Message)
	}

	resp := dto.UserResponse{
		ID:    user.ID,
		Login: user.Login,
		Name:  user.Name,
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	err := h.service.DeleteUser(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Message)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	var req dto.UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := &models.UserUpdate{
		Name: req.Name,
	}
	err := h.service.UpdateUser(ctx, userID, user)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Message)
	}

	return c.NoContent(http.StatusOK)
}
