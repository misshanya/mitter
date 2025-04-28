package mitt

import (
	"context"
	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/models"
)

type Service struct {
	mr models.MittRepository
}

func NewService(mr models.MittRepository) *Service {
	return &Service{mr: mr}
}

func (s *Service) CreateMitt(ctx context.Context, userID uuid.UUID, mitt *models.MittCreate) (*models.Mitt, error) {
	return s.mr.CreateMitt(ctx, userID, mitt)
}

func (s *Service) GetMitt(ctx context.Context, id uuid.UUID) (*models.Mitt, error) {
	return s.mr.GetMitt(ctx, id)
}

func (s *Service) GetAllUserMitts(ctx context.Context, userID uuid.UUID) ([]*models.Mitt, error) {
	return s.mr.GetAllUserMitts(ctx, userID)
}

func (s *Service) UpdateMitt(ctx context.Context, mittID uuid.UUID, mitt *models.MittUpdate) (*models.Mitt, error) {
	return s.mr.UpdateMitt(ctx, mittID, mitt)
}

func (s *Service) DeleteMitt(ctx context.Context, mittID uuid.UUID) error {
	return s.mr.DeleteMitt(ctx, mittID)
}
