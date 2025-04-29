package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/db/sqlc/storage"
	"github.com/misshanya/mitter/internal/models"
)

type MittRepository struct {
	queries *storage.Queries
}

func NewMittRepository(q *storage.Queries) *MittRepository {
	return &MittRepository{queries: q}
}

func mittDBToMitt(mittDB storage.Mitt) *models.Mitt {
	return &models.Mitt{
		ID:        mittDB.ID,
		Author:    mittDB.Author,
		Content:   mittDB.Content,
		CreatedAt: mittDB.CreatedAt.Time,
		UpdatedAt: mittDB.UpdatedAt.Time,
	}
}

func (r *MittRepository) CreateMitt(ctx context.Context, userID uuid.UUID, mitt *models.MittCreate) (*models.Mitt, error) {
	mittDB, err := r.queries.CreateMitt(ctx, storage.CreateMittParams{
		Author:  userID,
		Content: mitt.Content,
	})
	if err != nil {
		return nil, err
	}

	return mittDBToMitt(mittDB), nil
}

func (r *MittRepository) GetMitt(ctx context.Context, id uuid.UUID) (*models.Mitt, error) {
	mittDB, err := r.queries.GetMitt(ctx, id)
	if err != nil {
		return nil, err
	}

	return mittDBToMitt(mittDB), nil
}

func (r *MittRepository) GetAllUserMitts(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*models.Mitt, error) {
	mittsDB, err := r.queries.GetAllUserMitts(ctx, storage.GetAllUserMittsParams{
		Author: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	mitts := make([]*models.Mitt, len(mittsDB))
	for i, mittDB := range mittsDB {
		mitts[i] = mittDBToMitt(mittDB)
	}

	return mitts, nil
}

func (r *MittRepository) UpdateMitt(ctx context.Context, mittID uuid.UUID, mitt *models.MittUpdate) (*models.Mitt, error) {
	mittDB, err := r.queries.UpdateMitt(ctx, storage.UpdateMittParams{
		ID:      mittID,
		Content: mitt.Content,
	})
	if err != nil {
		return nil, err
	}

	return mittDBToMitt(mittDB), nil
}

func (r *MittRepository) DeleteMitt(ctx context.Context, mittID uuid.UUID) error {
	return r.queries.DeleteMitt(ctx, mittID)
}
