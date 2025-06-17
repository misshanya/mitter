package mitt

import (
	"context"
	"errors"
	"github.com/misshanya/mitter/pkg/pgutil"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/misshanya/mitter/internal/models"
)

type Service struct {
	mr models.MittRepository
	mm models.MittMetrics
	ur models.UserRepository
	l  *slog.Logger
}

func NewService(mr models.MittRepository, mm models.MittMetrics, ur models.UserRepository, l *slog.Logger) *Service {
	return &Service{mr: mr, mm: mm, ur: ur, l: l}
}

func (s *Service) CreateMitt(ctx context.Context, userID uuid.UUID, mitt *models.MittCreate) (*models.Mitt, *models.HTTPError) {
	newMitt, err := s.mr.CreateMitt(ctx, userID, mitt)
	if err != nil {
		s.l.Error("error creating mitt", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	// Update metrics
	go s.mm.AddMitt()

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
		s.l.Error("error getting mitt", slog.Any("err", err))
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
		s.l.Error("error getting mitts", slog.Any("err", err))
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
	if existingMitt.AuthorID != userID {
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
		s.l.Error("error updating mitt", slog.Any("err", err))
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
	if existingMitt.AuthorID != userID {
		return &models.HTTPError{
			Code:    http.StatusForbidden,
			Message: "You are not allowed to do this",
		}
	}

	err := s.mr.DeleteMitt(ctx, mittID)
	if err != nil {
		s.l.Error("error deleting mitt", slog.Any("err", err))
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	// Update metrics
	go s.mm.DeleteMitt()

	return nil
}

// Likes

func (s *Service) addLike(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) error {
	err := s.mr.LikeMitt(ctx, userID, mittID)
	if err != nil {
		return err
	}

	// Add like in metrics
	go s.mm.AddLike()

	return nil
}

func (s *Service) deleteLike(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) error {
	err := s.mr.DeleteMittLike(ctx, userID, mittID)
	if err != nil {
		return err
	}

	// Delete like in metrics
	go s.mm.DeleteLike()

	return nil
}

func (s *Service) SwitchLike(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) (bool, *models.HTTPError) {
	err := s.addLike(ctx, userID, mittID)
	if err != nil && !pgutil.IsUniqueViolation(err) {
		// If mitt doesn't exist
		if pgutil.IsForeignKeyViolation(err) {
			return false, &models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Mitt doesn't exist",
			}
		}

		s.l.Error("error liking mitt", slog.Any("err", err))
		return false, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	} else if pgutil.IsUniqueViolation(err) {
		// Delete like
		err := s.deleteLike(ctx, userID, mittID)
		if err != nil {
			s.l.Error("error deleting mitt like", slog.Any("err", err))
			return false, &models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
		}

		return false, nil
	}

	return true, nil
}

func (s *Service) Feed(ctx context.Context, limit, offset int32) ([]*models.Mitt, *models.HTTPError) {
	mitts, err := s.mr.Feed(ctx, limit, offset)
	if err != nil {
		s.l.Error("failed to get feed", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	go s.mm.ViewInFeed(float64(len(mitts)))

	return mitts, nil
}
