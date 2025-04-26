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
	CreateUser(ctx context.Context, user *models.UserCreate) (uuid.UUID, *models.HTTPError)
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
	group.POST("/", h.createUser)
}

func (h *UserHandler) createUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.UserCreateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := models.UserCreate{
		Login: req.Login,
		Name:  req.Name,
	}

	id, err := h.service.CreateUser(ctx, &user)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Message)
	}

	resp := dto.UserCreateResponse{
		ID: id,
	}
	return c.JSON(http.StatusCreated, resp)
}
