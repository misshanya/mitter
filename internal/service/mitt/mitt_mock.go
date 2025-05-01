package mitt

import (
	"context"
	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/models"
	"time"
)

var mockMittModel = &models.Mitt{
	ID:        uuid.New(),
	Author:    uuid.New(),
	Content:   "hello world",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Likes:     0,
}

// Mock mitt repo
type mockMittRepo struct{}

func (m mockMittRepo) CreateMitt(ctx context.Context, userID uuid.UUID, mitt *models.MittCreate) (*models.Mitt, error) {
	_ = ctx
	_ = userID
	_ = mitt

	return mockMittModel, nil
}

func (m mockMittRepo) GetMitt(ctx context.Context, id uuid.UUID) (*models.Mitt, error) {
	_ = ctx
	_ = id

	return mockMittModel, nil
}

func (m mockMittRepo) GetAllUserMitts(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*models.Mitt, error) {
	_ = ctx
	_ = userID
	_ = limit
	_ = offset

	return []*models.Mitt{mockMittModel}, nil
}

func (m mockMittRepo) UpdateMitt(ctx context.Context, mittID uuid.UUID, mitt *models.MittUpdate) (*models.Mitt, error) {
	_ = ctx
	_ = mittID

	// Update mock mitt content
	mockMittModel.Content = mitt.Content

	return mockMittModel, nil
}

func (m mockMittRepo) DeleteMitt(ctx context.Context, mittID uuid.UUID) error {
	_ = ctx
	_ = mittID

	// Nothing happens :)
	return nil
}

func (m mockMittRepo) LikeMitt(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) error {
	_ = ctx
	_ = userID
	_ = mittID

	// Add like to mock mitt
	mockMittModel.Likes++

	return nil
}

func (m mockMittRepo) IsMittLikedByUser(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) (bool, error) {
	_ = ctx
	_ = userID
	_ = mittID

	// Let's assume that's true.
	return true, nil
}

func (m mockMittRepo) DeleteMittLike(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) error {
	_ = ctx
	_ = userID
	_ = mittID

	// Delete like (oops, it can be negative cause you can run tests in different order)
	mockMittModel.Likes--

	return nil
}

func (m mockMittRepo) GetMittLikesCount(ctx context.Context, mittID uuid.UUID) (int64, error) {
	_ = ctx
	_ = mittID

	return mockMittModel.Likes, nil
}
