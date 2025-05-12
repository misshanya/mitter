package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

// Likes

func (r *MittRepository) LikeMitt(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) error {
	return r.queries.LikeMitt(ctx, storage.LikeMittParams{
		UserID: userID,
		MittID: mittID,
	})
}

func (r *MittRepository) IsMittLikedByUser(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) (bool, error) {
	_, err := r.queries.IsMittLikedByUser(ctx, storage.IsMittLikedByUserParams{
		UserID: userID,
		MittID: mittID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *MittRepository) DeleteMittLike(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) error {
	return r.queries.DeleteMittLike(ctx, storage.DeleteMittLikeParams{
		UserID: userID,
		MittID: mittID,
	})
}

func (r *MittRepository) GetMittLikesCount(ctx context.Context, mittID uuid.UUID) (int64, error) {
	return r.queries.GetMittLikesCount(ctx, mittID)
}

func (r *MittRepository) Feed(ctx context.Context, limit, offset int32) ([]*models.Mitt, error) {
	mittsDB, err := r.queries.Feed(ctx, storage.FeedParams{
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
