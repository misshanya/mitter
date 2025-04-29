package mitt

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/misshanya/mitter/internal/models"
	"log/slog"
	"net/http"
)

type Service struct {
	mr models.MittRepository
}

func NewService(mr models.MittRepository) *Service {
	return &Service{mr: mr}
}

func (s *Service) CreateMitt(ctx context.Context, userID uuid.UUID, mitt *models.MittCreate) (*models.Mitt, *models.HTTPError) {
	newMitt, err := s.mr.CreateMitt(ctx, userID, mitt)
	if err != nil {
		slog.Error("error creating mitt", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	return newMitt, nil
}

func (s *Service) GetMitt(ctx context.Context, id uuid.UUID) (*models.Mitt, *models.HTTPError) {
	mitt, err := s.mr.GetMitt(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Mitt not found",
			}
		}
		slog.Error("error getting mitt", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	return mitt, nil
}

func (s *Service) GetAllUserMitts(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*models.Mitt, *models.HTTPError) {
	mitts, err := s.mr.GetAllUserMitts(ctx, userID, limit, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Mitts not found",
			}
		}
		slog.Error("error getting mitts", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	return mitts, nil
}

func (s *Service) UpdateMitt(ctx context.Context, userID uuid.UUID, mittID uuid.UUID, mitt *models.MittUpdate) (*models.Mitt, *models.HTTPError) {
	existingMitt, httpErr := s.GetMitt(ctx, mittID)
	if httpErr != nil {
		return nil, httpErr
	}

	// Check if user is author of mitt
	if existingMitt.Author != userID {
		return nil, &models.HTTPError{
			Code:    http.StatusForbidden,
			Message: "You are not allowed to do this",
		}
	}

	newMitt, err := s.mr.UpdateMitt(ctx, mittID, mitt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Mitt not found",
			}
		}
		slog.Error("error updating mitt", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	return newMitt, nil
}

func (s *Service) DeleteMitt(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) *models.HTTPError {
	existingMitt, httpErr := s.GetMitt(ctx, mittID)
	if httpErr != nil {
		return httpErr
	}

	// Check if user is author of mitt
	if existingMitt.Author != userID {
		return &models.HTTPError{
			Code:    http.StatusForbidden,
			Message: "You are not allowed to do this",
		}
	}

	err := s.mr.DeleteMitt(ctx, mittID)
	if err != nil {
		slog.Error("error deleting mitt", slog.Any("err", err))
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	return nil
}
