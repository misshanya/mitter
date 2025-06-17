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

func (r *MittRepository) CreateMitt(ctx context.Context, userID uuid.UUID, mitt *models.MittCreate) (*models.Mitt, error) {
	row, err := r.queries.CreateMitt(ctx, storage.CreateMittParams{
		Author:  userID,
		Content: mitt.Content,
	})
	if err != nil {
		return nil, err
	}

	newMitt := &models.Mitt{
		ID:         row.ID,
		AuthorID:   row.Author,
		AuthorName: row.AuthorName.String,
		Content:    row.Content,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
		Likes:      row.LikesCount,
	}

	return newMitt, nil
}

func (r *MittRepository) GetMitt(ctx context.Context, id uuid.UUID) (*models.Mitt, error) {
	row, err := r.queries.GetMitt(ctx, id)
	if err != nil {
		return nil, err
	}

	newMitt := &models.Mitt{
		ID:         row.ID,
		AuthorID:   row.Author,
		AuthorName: row.AuthorName.String,
		Content:    row.Content,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
		Likes:      row.LikesCount,
	}

	return newMitt, nil
}

func (r *MittRepository) GetAllUserMitts(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*models.Mitt, error) {
	rows, err := r.queries.GetAllUserMitts(ctx, storage.GetAllUserMittsParams{
		Author: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	mitts := make([]*models.Mitt, len(rows))
	for i, row := range rows {
		mitts[i] = &models.Mitt{
			ID:         row.ID,
			AuthorID:   row.Author,
			AuthorName: row.AuthorName.String,
			Content:    row.Content,
			CreatedAt:  row.CreatedAt.Time,
			UpdatedAt:  row.UpdatedAt.Time,
			Likes:      row.LikesCount,
		}
	}

	return mitts, nil
}

func (r *MittRepository) UpdateMitt(ctx context.Context, mittID uuid.UUID, mitt *models.MittUpdate) (*models.Mitt, error) {
	row, err := r.queries.UpdateMitt(ctx, storage.UpdateMittParams{
		ID:      mittID,
		Content: mitt.Content,
	})
	if err != nil {
		return nil, err
	}

	newMitt := &models.Mitt{
		ID:         row.ID,
		AuthorID:   row.Author,
		AuthorName: row.AuthorName.String,
		Content:    row.Content,
		CreatedAt:  row.CreatedAt.Time,
		UpdatedAt:  row.UpdatedAt.Time,
		Likes:      row.LikesCount,
	}

	return newMitt, nil
}

func (r *MittRepository) DeleteMitt(ctx context.Context, mittID uuid.UUID) error {
	return r.queries.DeleteMitt(ctx, mittID)
}

// Likes

func (r *MittRepository) LikeMitt(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) error {
	return r.queries.LikeMitt(ctx, storage.LikeMittParams{
		UserID: userID,
		MittID: mittID,
	})
}

func (r *MittRepository) DeleteMittLike(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) error {
	return r.queries.DeleteMittLike(ctx, storage.DeleteMittLikeParams{
		UserID: userID,
		MittID: mittID,
	})
}

func (r *MittRepository) Feed(ctx context.Context, limit, offset int32) ([]*models.Mitt, error) {
	rows, err := r.queries.Feed(ctx, storage.FeedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	mitts := make([]*models.Mitt, len(rows))
	for i, row := range rows {
		mitts[i] = &models.Mitt{
			ID:         row.ID,
			AuthorID:   row.Author,
			AuthorName: row.AuthorName.String,
			Content:    row.Content,
			CreatedAt:  row.CreatedAt.Time,
			UpdatedAt:  row.UpdatedAt.Time,
			Likes:      row.LikesCount,
		}
	}

	return mitts, nil
}
