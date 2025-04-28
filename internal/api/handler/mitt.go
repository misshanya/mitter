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

type mittService interface {
	CreateMitt(ctx context.Context, userID uuid.UUID, mitt *models.MittCreate) (*models.Mitt, *models.HTTPError)

	GetMitt(ctx context.Context, id uuid.UUID) (*models.Mitt, *models.HTTPError)
	GetAllUserMitts(ctx context.Context, userID uuid.UUID) ([]*models.Mitt, *models.HTTPError)

	UpdateMitt(ctx context.Context, userID uuid.UUID, mittID uuid.UUID, mitt *models.MittUpdate) (*models.Mitt, *models.HTTPError)

	DeleteMitt(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) *models.HTTPError
}

type MittHandler struct {
	ms                mittService
	validate          *validator.Validate
	reqAuthMiddleware echo.MiddlewareFunc
}

func NewMittHandler(ms mittService, reqAuthMdl echo.MiddlewareFunc) *MittHandler {
	return &MittHandler{
		ms:                ms,
		validate:          validator.New(),
		reqAuthMiddleware: reqAuthMdl,
	}
}

func (h *MittHandler) Routes(group *echo.Group) {
	group.POST("", h.createMitt, h.reqAuthMiddleware)
	group.GET("/:id", h.getMitt)
	group.GET("/user/:id", h.getAllUserMitts)
	group.PUT("/:id", h.updateMitt, h.reqAuthMiddleware)
	group.DELETE("/:id", h.deleteMitt, h.reqAuthMiddleware)
}

// createMitt godoc
//
//	@Summary	Create Mitt
//	@Tags		Mitts
//	@Security	Bearer
//	@Param		Authorization	header	string	true	"access token 'Bearer {token}'"
//	@Accept		json
//	@Param		CreateMittRequest	body	dto.MittCreateRequest	true	"Create Mitt Request"
//	@Produce	json
//	@Success	201	{object}	dto.MittResponse
//	@Failure	400	{object}	dto.HTTPError
//	@Failure	401	{object}	dto.HTTPError
//	@Failure	500	{object}	dto.HTTPError
//	@Router		/mitt [post]
func (h *MittHandler) createMitt(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	var req dto.MittCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	// Validate
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	mittCreate := &models.MittCreate{
		Content: req.Content,
	}
	mitt, err := h.ms.CreateMitt(ctx, userID, mittCreate)
	if err != nil {
		return c.JSON(err.Code, dto.HTTPError{Message: err.Message})
	}

	resp := dto.MittResponse{
		ID:        mitt.ID,
		Author:    mitt.Author,
		Content:   mitt.Content,
		CreatedAt: mitt.CreatedAt,
		UpdatedAt: mitt.UpdatedAt,
	}
	return c.JSON(http.StatusCreated, resp)
}

// getMitt godoc
//
//	@Summary	Get Mitt
//	@Tags		Mitts
//	@Param		id	path	string	true	"ID of mitt"
//	@Produce	json
//	@Success	200	{object}	dto.MittResponse
//	@Failure	400	{object}	dto.HTTPError
//	@Failure	404	{object}	dto.HTTPError
//	@Failure	500	{object}	dto.HTTPError
//	@Router		/mitt/{id} [get]
func (h *MittHandler) getMitt(c echo.Context) error {
	ctx := c.Request().Context()

	mittIDStr := c.Param("id")
	mittID, err := uuid.Parse(mittIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	mitt, httpErr := h.ms.GetMitt(ctx, mittID)
	if httpErr != nil {
		return c.JSON(httpErr.Code, dto.HTTPError{Message: httpErr.Message})
	}

	resp := dto.MittResponse{
		ID:        mitt.ID,
		Author:    mitt.Author,
		Content:   mitt.Content,
		CreatedAt: mitt.CreatedAt,
		UpdatedAt: mitt.UpdatedAt,
	}
	return c.JSON(http.StatusOK, resp)
}

// getAllUserMitts godoc
//
//	@Summary	Get User Mitts
//	@Tags		Mitts
//	@Param		id	path	string	true	"ID of user"
//	@Produce	json
//	@Success	200	{object}	[]dto.MittResponse
//	@Failure	400	{object}	dto.HTTPError
//	@Failure	404	{object}	dto.HTTPError
//	@Failure	500	{object}	dto.HTTPError
//	@Router		/mitt/user/{id} [get]
func (h *MittHandler) getAllUserMitts(c echo.Context) error {
	ctx := c.Request().Context()

	userIDToGetStr := c.Param("id")
	userIDToGet, err := uuid.Parse(userIDToGetStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	mitts, httpErr := h.ms.GetAllUserMitts(ctx, userIDToGet)
	if httpErr != nil {
		return c.JSON(httpErr.Code, dto.HTTPError{Message: httpErr.Message})
	}

	resp := make([]dto.MittResponse, len(mitts))
	for i, m := range mitts {
		resp[i] = dto.MittResponse{
			ID:        m.ID,
			Author:    m.Author,
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	return c.JSON(http.StatusOK, resp)
}

// updateMitt godoc
//
//	@Summary	Update Mitt
//	@Tags		Mitts
//	@Security	Bearer
//	@Param		Authorization	header	string	true	"access token 'Bearer {token}'"
//	@Param		id				path	string	true	"ID of mitt"
//	@Accept		json
//	@Param		UpdateMittRequest	body	dto.MittUpdateRequest	true	"Update Mitt Request"
//	@Produce	json
//	@Success	200	{object}	dto.MittResponse
//	@Failure	400	{object}	dto.HTTPError
//	@Failure	401	{object}	dto.HTTPError
//	@Failure	403	{object}	dto.HTTPError
//	@Failure	404	{object}	dto.HTTPError
//	@Failure	500	{object}	dto.HTTPError
//	@Router		/mitt/{id} [put]
func (h *MittHandler) updateMitt(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	mittIDStr := c.Param("id")
	mittID, err := uuid.Parse(mittIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	var req dto.MittUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	// Validate
	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	updateMitt := &models.MittUpdate{
		Content: req.Content,
	}
	newMitt, httpErr := h.ms.UpdateMitt(ctx, userID, mittID, updateMitt)
	if httpErr != nil {
		return c.JSON(httpErr.Code, dto.HTTPError{Message: httpErr.Message})
	}
	return c.JSON(http.StatusOK, newMitt)
}

// deleteMitt godoc
//
//	@Tags		Mitts
//	@Summary	Delete mitt
//	@Security	Bearer
//	@Param		Authorization	header	string	true	"access token 'Bearer {token}'"
//	@Param		id				path	string	true	"ID of mitt"
//	@Produce	json
//	@Success	204
//	@Failure	400	{object}	dto.HTTPError
//	@Failure	401	{object}	dto.HTTPError
//	@Failure	403	{object}	dto.HTTPError
//	@Failure	500	{object}	dto.HTTPError
//	@Router		/mitt/{id} [delete]
func (h *MittHandler) deleteMitt(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Get("userID").(uuid.UUID)

	mittIDStr := c.Param("id")
	mittID, err := uuid.Parse(mittIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.HTTPError{Message: err.Error()})
	}

	if httpErr := h.ms.DeleteMitt(ctx, userID, mittID); httpErr != nil {
		return c.JSON(httpErr.Code, dto.HTTPError{Message: httpErr.Message})
	}
	return c.NoContent(http.StatusNoContent)
}
