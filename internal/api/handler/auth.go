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

	ChangePassword(ctx context.Context, id uuid.UUID, changePassword *models.ChangePassword) *models.HTTPError
}

type AuthHandler struct {
	as                authService
	validate          *validator.Validate
	reqAuthMiddleware echo.MiddlewareFunc
}

func NewAuthHandler(ar authService, reqAuthMdl echo.MiddlewareFunc) *AuthHandler {
	return &AuthHandler{
		as:                ar,
		validate:          validator.New(),
		reqAuthMiddleware: reqAuthMdl,
	}
}

func (h *AuthHandler) Routes(group *echo.Group) {
	group.POST("/sign-in", h.signIn)
	group.POST("/sign-up", h.signUp)

	// Protect /change-password with auth middleware
	group.POST("/change-password", h.changePassword, h.reqAuthMiddleware)
}

// signIn godoc
//
//	@Summary		Sign In
//	@Description	Sign In user via login and password
//	@Tags			Auth
//	@Accept			json
//	@Param			SignInRequest	body	dto.SignInRequest	true	"Sign In Request"
//	@Produce		json
//	@Success		200	{object}	dto.SignInResponse
//	@Failure		400	{object}	dto.HTTPError
//	@Failure		401	{object}	dto.HTTPError
//	@Failure		500	{object}	dto.HTTPError
//	@Router			/auth/sign-in [post]
func (h *AuthHandler) signIn(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.SignInRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	// Validate
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	creds := models.SignIn{
		Login:    req.Login,
		Password: req.Password,
	}

	token, err := h.as.SignIn(ctx, creds)
	if err != nil {
		return c.JSON(err.Code, dto.HTTPError{Message: err.Message})
	}

	resp := dto.SignInResponse{
		Token: token,
	}
	return c.JSON(http.StatusOK, resp)
}

// signUp godoc
//
//	@Summary		Sign Up
//	@Description	Sign Up user
//	@Tags			Auth
//	@Accept			json
//	@Param			SignUpRequest	body	dto.SignUpRequest	true	"Sign Up Request"
//	@Produce		json
//	@Success		201	{object}	dto.SignUpResponse
//	@Failure		400	{object}	dto.HTTPError
//	@Failure		401	{object}	dto.HTTPError
//	@Failure		409	{object}	dto.HTTPError
//	@Failure		500	{object}	dto.HTTPError
//	@Router			/auth/sign-up [post]
func (h *AuthHandler) signUp(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	// Validate
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	user := &models.UserCreate{
		Login:    req.Login,
		Name:     req.Name,
		Password: req.Password,
	}
	id, err := h.as.SignUp(ctx, user)
	if err != nil {
		return c.JSON(err.Code, dto.HTTPError{Message: err.Message})
	}

	resp := dto.SignUpResponse{
		ID: id,
	}
	return c.JSON(http.StatusCreated, resp)
}

// changePassword godoc
//
//	@Summary		Change Password
//	@Description	Change user's password
//	@Tags			Auth
//	@Security		Bearer
//	@Param			Authorization	header	string	true	"access token 'Bearer {token}'"
//	@Accept			json
//	@Param			ChangePasswordRequest	body	dto.ChangePasswordRequest	true	"Change Password Request"
//	@Success		200
//	@Failure		400	{object}	dto.HTTPError
//	@Failure		401	{object}	dto.HTTPError
//	@Failure		500	{object}	dto.HTTPError
//	@Router			/auth/change-password [post]
func (h *AuthHandler) changePassword(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	// Validate
	if err := h.validate.StructCtx(ctx, &req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	changePwd := &models.ChangePassword{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
	err := h.as.ChangePassword(ctx, userID, changePwd)
	if err != nil {
		return c.JSON(err.Code, dto.HTTPError{Message: err.Message})
	}

	return c.NoContent(http.StatusOK)
}
