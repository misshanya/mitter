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
	group.GET("", h.GetMe)
	group.DELETE("", h.DeleteUser)
	group.PATCH("", h.UpdateUser)
}

// GetMe godoc
//
//	@Tags			User
//	@Summary		Get Me
//	@Description	Get info about me
//	@Security		Bearer
//	@Param			Authorization	header	string	true	"access token 'Bearer {token}'"
//	@Produce		json
//	@Success		200	{object}	dto.UserResponse
//	@Failure		400	{object}	dto.HTTPError
//	@Failure		401	{object}	dto.HTTPError
//	@Failure		404	{object}	dto.HTTPError
//	@Failure		500	{object}	dto.HTTPError
//	@Router			/user [get]
func (h *UserHandler) GetMe(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		return c.JSON(err.Code, dto.HTTPError{Message: err.Message})
	}

	resp := dto.UserResponse{
		ID:    user.ID,
		Login: user.Login,
		Name:  user.Name,
	}
	return c.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
//
//	@Tags			User
//	@Summary		Delete account
//	@Description	Delete user account
//	@Security		Bearer
//	@Param			Authorization	header	string	true	"access token 'Bearer {token}'"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	dto.HTTPError
//	@Failure		401	{object}	dto.HTTPError
//	@Failure		500	{object}	dto.HTTPError
//	@Router			/user [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	err := h.service.DeleteUser(ctx, userID)
	if err != nil {
		return c.JSON(err.Code, dto.HTTPError{Message: err.Message})
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateUser godoc
//
//	@Tags			User
//	@Summary		Update user
//	@Description	Update user info
//	@Security		Bearer
//	@Param			Authorization		header	string					true	"access token 'Bearer {token}'"
//	@Param			UpdateUserRequest	body	dto.UserUpdateRequest	true	"Update User Request"
//	@Accept			json
//	@Success		200
//	@Failure		400	{object}	dto.HTTPError
//	@Failure		401	{object}	dto.HTTPError
//	@Failure		500	{object}	dto.HTTPError
//	@Router			/user [patch]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	var req dto.UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	// Validate
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	user := &models.UserUpdate{
		Name: req.Name,
	}
	err := h.service.UpdateUser(ctx, userID, user)
	if err != nil {
		return c.JSON(err.Code, dto.HTTPError{Message: err.Message})
	}

	return c.NoContent(http.StatusOK)
}
