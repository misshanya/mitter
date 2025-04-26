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

type authService interface {
	SignIn(ctx context.Context, creds models.SignIn) (string, *models.HTTPError)
	SignUp(ctx context.Context, user *models.UserCreate) (uuid.UUID, *models.HTTPError)
}

type AuthHandler struct {
	as       authService
	validate *validator.Validate
}

func NewAuthHandler(ar authService) *AuthHandler {
	return &AuthHandler{
		as:       ar,
		validate: validator.New(),
	}
}

func (h *AuthHandler) Routes(group *echo.Group) {
	group.POST("/sign-in", h.signIn)
	group.POST("/sign-up", h.signUp)
}

func (h *AuthHandler) signIn(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.SignInRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	creds := models.SignIn{
		Login:    req.Login,
		Password: req.Password,
	}

	token, err := h.as.SignIn(ctx, creds)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Message)
	}

	resp := dto.SignInResponse{
		Token: token,
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) signUp(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.SignUpRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := &models.UserCreate{
		Login:    req.Login,
		Name:     req.Name,
		Password: req.Password,
	}
	id, err := h.as.SignUp(ctx, user)
	if err != nil {
		return echo.NewHTTPError(err.Code, err.Message)
	}

	resp := dto.SignUpResponse{
		ID: id,
	}
	return c.JSON(http.StatusCreated, resp)
}
