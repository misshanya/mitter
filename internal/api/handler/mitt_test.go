package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/misshanya/mitter/internal/api/dto"
	"github.com/misshanya/mitter/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mockMittModel = &models.Mitt{
	ID:        uuid.New(),
	AuthorID:  mockUserID,
	Content:   "hello world",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Likes:     0,
}

var mockUserID = uuid.MustParse("b096376a-5fa9-4130-907a-709c67008a65")

// Mock service
type mockMittService struct{}

func (s *mockMittService) CreateMitt(ctx context.Context, userID uuid.UUID, mitt *models.MittCreate) (*models.Mitt, *models.HTTPError) {
	_ = ctx
	_ = userID
	_ = mitt

	return mockMittModel, nil
}

func (s *mockMittService) GetMitt(ctx context.Context, id uuid.UUID) (*models.Mitt, *models.HTTPError) {
	_ = ctx
	_ = id

	return mockMittModel, nil
}

func (s *mockMittService) GetAllUserMitts(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*models.Mitt, *models.HTTPError) {
	_ = ctx
	_ = userID
	_ = limit
	_ = offset

	return []*models.Mitt{mockMittModel}, nil
}

func (m *mockMittService) UpdateMitt(ctx context.Context, userID uuid.UUID, mittID uuid.UUID, mitt *models.MittUpdate) (*models.Mitt, *models.HTTPError) {
	_ = ctx
	_ = mittID

	// Update mock mitt content
	mockMittModel.Content = mitt.Content

	return mockMittModel, nil
}

func (m *mockMittService) DeleteMitt(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) *models.HTTPError {
	_ = ctx
	_ = userID
	_ = mittID

	// Do nothing
	return nil
}

func (m *mockMittService) SwitchLike(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) (bool, *models.HTTPError) {
	// If likes > 0, lets suppose user already liked mitt
	isAlreadyLiked := mockMittModel.Likes > 0

	if !isAlreadyLiked {
		mockMittModel.Likes++
		return true, nil
	}

	mockMittModel.Likes--

	return false, nil
}

func (m *mockMittService) Feed(ctx context.Context, limit, offset int32) ([]*models.Mitt, *models.HTTPError) {
	_ = ctx
	_ = limit
	_ = offset

	return []*models.Mitt{mockMittModel}, nil
}

// Tests
func TestMittHandler_CreateMitt(t *testing.T) {
	e := echo.New()
	handler := NewMittHandler(&mockMittService{}, mockRequireAuth)

	g := e.Group("/api/v1/mitt")
	handler.Routes(g)

	// Create request
	reqBody := `{"content":"hello world"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/mitt", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	if assert.NoError(t, mockRequireAuth(handler.createMitt)(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		resp := dto.MittResponse{
			ID:        mockMittModel.ID,
			Author:    mockMittModel.AuthorID,
			Content:   mockMittModel.Content,
			CreatedAt: mockMittModel.CreatedAt,
			UpdatedAt: mockMittModel.UpdatedAt,
			Likes:     mockMittModel.Likes,
		}

		b, err := json.Marshal(resp)
		require.NoError(t, err)

		expectedResp := string(b)

		require.JSONEq(t, expectedResp, rec.Body.String())
	}
}

func TestMittHandler_GetMitt(t *testing.T) {
	e := echo.New()
	handler := NewMittHandler(&mockMittService{}, mockRequireAuth)

	g := e.Group("/api/v1/mitt")
	handler.Routes(g)

	// Create request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/mitt/%s", mockMittModel.ID.String()), nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	// Set path param (id)
	ctx.SetPath("/api/v1/mitt/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(mockMittModel.ID.String())

	if assert.NoError(t, handler.getMitt(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := dto.MittResponse{
			ID:        mockMittModel.ID,
			Author:    mockMittModel.AuthorID,
			Content:   mockMittModel.Content,
			CreatedAt: mockMittModel.CreatedAt,
			UpdatedAt: mockMittModel.UpdatedAt,
			Likes:     mockMittModel.Likes,
		}

		b, err := json.Marshal(resp)
		require.NoError(t, err)

		expectedResp := string(b)

		require.JSONEq(t, expectedResp, rec.Body.String())
	}
}

func TestMittHandler_GetAllUserMitts(t *testing.T) {
	e := echo.New()
	handler := NewMittHandler(&mockMittService{}, mockRequireAuth)

	g := e.Group("/api/v1/mitt")
	handler.Routes(g)

	// Create request
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/mitt/user/%s", mockMittModel.AuthorID.String()), nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	// Set path param (id)
	ctx.SetPath("/api/v1/mitt/user/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(mockUserID.String())

	if assert.NoError(t, handler.getAllUserMitts(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		respMitt := dto.MittResponse{
			ID:        mockMittModel.ID,
			Author:    mockMittModel.AuthorID,
			Content:   mockMittModel.Content,
			CreatedAt: mockMittModel.CreatedAt,
			UpdatedAt: mockMittModel.UpdatedAt,
			Likes:     mockMittModel.Likes,
		}

		resp := []dto.MittResponse{respMitt}

		b, err := json.Marshal(resp)
		require.NoError(t, err)

		expectedResp := string(b)

		require.JSONEq(t, expectedResp, rec.Body.String())
	}
}

func TestMittHandler_UpdateMitt(t *testing.T) {
	e := echo.New()
	handler := NewMittHandler(&mockMittService{}, mockRequireAuth)

	g := e.Group("/api/v1/mitt")
	handler.Routes(g)

	// Create request
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/mitt/%s", mockMittModel.ID.String()), nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	// Set path param (id)
	ctx.SetPath("/api/v1/mitt/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues(mockMittModel.ID.String())

	if assert.NoError(t, mockRequireAuth(handler.updateMitt)(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := dto.MittResponse{
			ID:        mockMittModel.ID,
			Author:    mockMittModel.AuthorID,
			Content:   mockMittModel.Content,
			CreatedAt: mockMittModel.CreatedAt,
			UpdatedAt: mockMittModel.UpdatedAt,
			Likes:     mockMittModel.Likes,
		}

		b, err := json.Marshal(resp)
		require.NoError(t, err)

		expectedResp := string(b)

		require.JSONEq(t, expectedResp, rec.Body.String())
	}
}

func TestMittHandler_DeleteMitt(t *testing.T) {
	e := echo.New()
	handler := NewMittHandler(&mockMittService{}, mockRequireAuth)

	g := e.Group("/api/v1/mitt")
	handler.Routes(g)

	// Create request
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/mitt/%s/like", mockMittModel.ID.String()), nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	// Set path param (id)
	ctx.SetPath("/api/v1/mitt/:id/like")
	ctx.SetParamNames("id")
	ctx.SetParamValues(mockMittModel.ID.String())

	if assert.NoError(t, mockRequireAuth(handler.deleteMitt)(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func TestMittHandler_LikeMitt(t *testing.T) {
	e := echo.New()
	handler := NewMittHandler(&mockMittService{}, mockRequireAuth)

	g := e.Group("/api/v1/mitt")
	handler.Routes(g)

	// Create request
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/mitt/%s/like", mockMittModel.ID.String()), nil)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	// Set path param (id)
	ctx.SetPath("/api/v1/mitt/:id/like")
	ctx.SetParamNames("id")
	ctx.SetParamValues(mockMittModel.ID.String())

	// Test like creation
	if assert.NoError(t, mockRequireAuth(handler.likeMitt)(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := dto.MittLikeResponse{
			Like: true,
		}

		b, err := json.Marshal(resp)
		require.NoError(t, err)

		expectedResp := string(b)

		require.JSONEq(t, expectedResp, rec.Body.String())
	}

	// Test like deletion

	// Create request
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/mitt/%s/like", mockMittModel.ID.String()), nil)
	rec = httptest.NewRecorder()

	ctx = e.NewContext(req, rec)

	// Set path param (id)
	ctx.SetPath("/api/v1/mitt/:id/like")
	ctx.SetParamNames("id")
	ctx.SetParamValues(mockMittModel.ID.String())

	if assert.NoError(t, mockRequireAuth(handler.likeMitt)(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := dto.MittLikeResponse{
			Like: false,
		}

		b, err := json.Marshal(resp)
		require.NoError(t, err)

		expectedResp := string(b)

		require.JSONEq(t, expectedResp, rec.Body.String())
	}
}
