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

	FollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) *models.HTTPError
	UnfollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) *models.HTTPError
	GetUserFollows(ctx context.Context, followerID uuid.UUID) ([]*models.User, *models.HTTPError)
	GetUserFollowers(ctx context.Context, followeeID uuid.UUID) ([]*models.User, *models.HTTPError)
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
	group.GET("", h.getMe)
	group.DELETE("", h.deleteUser)
	group.PATCH("", h.updateUser)

	group.POST("/:id/follow", h.followUser)
	group.DELETE("/:id/follow", h.unfollowUser)
	group.GET("/follows", h.getMyFollows)
	group.GET("/followers", h.getMyFollowers)
}

// getMe godoc
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
func (h *UserHandler) getMe(c echo.Context) error {
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

// deleteUser godoc
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
func (h *UserHandler) deleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	err := h.service.DeleteUser(ctx, userID)
	if err != nil {
		return c.JSON(err.Code, dto.HTTPError{Message: err.Message})
	}

	return c.NoContent(http.StatusNoContent)
}

// updateUser godoc
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
func (h *UserHandler) updateUser(c echo.Context) error {
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

// followUser godoc
//
//	@Tags		User
//	@Summary	Follow user
//	@Security	Bearer
//	@Param		Authorization	header	string	true	"access token 'Bearer {token}'"
//	@Param		id				path	string	true	"ID of user to follow to"
//	@Success	200
//	@Failure	400	{object}	dto.HTTPError
//	@Failure	401	{object}	dto.HTTPError
//	@Failure	404	{object}	dto.HTTPError
//	@Failure	409	{object}	dto.HTTPError
//	@Failure	500	{object}	dto.HTTPError
//	@Router		/user/{id}/follow [post]
func (h *UserHandler) followUser(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	followeeIDStr := c.Param("id")
	followeeID, err := uuid.Parse(followeeIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	httpErr := h.service.FollowUser(ctx, userID, followeeID)
	if httpErr != nil {
		return c.JSON(httpErr.Code, dto.HTTPError{Message: httpErr.Message})
	}

	return c.NoContent(http.StatusOK)
}

// unfollowUser godoc
//
//	@Tags		User
//	@Summary	Unfollow user
//	@Security	Bearer
//	@Param		Authorization	header	string	true	"access token 'Bearer {token}'"
//	@Param		id				path	string	true	"ID of user to unfollow to"
//	@Success	204
//	@Failure	400	{object}	dto.HTTPError
//	@Failure	401	{object}	dto.HTTPError
//	@Failure	500	{object}	dto.HTTPError
//	@Router		/user/{id}/follow [delete]
func (h *UserHandler) unfollowUser(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	followeeIDStr := c.Param("id")
	followeeID, err := uuid.Parse(followeeIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	httpErr := h.service.UnfollowUser(ctx, userID, followeeID)
	if httpErr != nil {
		return c.JSON(httpErr.Code, dto.HTTPError{Message: httpErr.Message})
	}

	return c.NoContent(http.StatusNoContent)
}

// getMyFollows godoc
//
//	@Tags		User
//	@Summary	Get My Follows
//	@Security	Bearer
//	@Param		Authorization	header		string	true	"access token 'Bearer {token}'"
//	@Success	200				{object}	[]dto.UserResponse
//	@Failure	400				{object}	dto.HTTPError
//	@Failure	401				{object}	dto.HTTPError
//	@Failure	500				{object}	dto.HTTPError
//	@Router		/user/follows [get]
func (h *UserHandler) getMyFollows(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	users, httpErr := h.service.GetUserFollows(ctx, userID)
	if httpErr != nil {
		return c.JSON(httpErr.Code, dto.HTTPError{Message: httpErr.Message})
	}

	resp := make([]dto.UserResponse, len(users))
	for i, user := range users {
		resp[i] = dto.UserResponse{
			ID:    user.ID,
			Login: user.Login,
			Name:  user.Name,
		}
	}

	return c.JSON(http.StatusOK, resp)
}

// getMyFollowers godoc
//
//	@Tags		User
//	@Summary	Get My Followers
//	@Security	Bearer
//	@Param		Authorization	header		string	true	"access token 'Bearer {token}'"
//	@Success	200				{object}	[]dto.UserResponse
//	@Failure	400				{object}	dto.HTTPError
//	@Failure	401				{object}	dto.HTTPError
//	@Failure	500				{object}	dto.HTTPError
//	@Router		/user/followers [get]
func (h *UserHandler) getMyFollowers(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	users, httpErr := h.service.GetUserFollowers(ctx, userID)
	if httpErr != nil {
		return c.JSON(httpErr.Code, dto.HTTPError{Message: httpErr.Message})
	}

	resp := make([]dto.UserResponse, len(users))
	for i, user := range users {
		resp[i] = dto.UserResponse{
			ID:    user.ID,
			Login: user.Login,
			Name:  user.Name,
		}
	}

	return c.JSON(http.StatusOK, resp)
}
