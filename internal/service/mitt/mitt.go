package mitt

import (
	"context"
	"errors"
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
}

func NewService(mr models.MittRepository, mm models.MittMetrics, ur models.UserRepository) *Service {
	return &Service{mr: mr, mm: mm, ur: ur}
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

	if err := s.setAuthorName(ctx, newMitt); err != nil {
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	// Update metrics
	go s.mm.AddMitt()

	return newMitt, nil
}

func (s *Service) setLikesCount(ctx context.Context, mitt *models.Mitt) error {
	likesCount, err := s.mr.GetMittLikesCount(ctx, mitt.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		slog.Error("error getting likes count", slog.Any("err", err))
		return err
	}
	mitt.Likes = likesCount
	return nil
}

func (s *Service) setAuthorName(ctx context.Context, mitt *models.Mitt) error {
	user, err := s.ur.GetUserByID(ctx, mitt.AuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		slog.Error("error getting user id", slog.Any("err", err))
		return err
	}
	mitt.AuthorName = user.Name
	return nil
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

	if err := s.setLikesCount(ctx, mitt); err != nil {
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	if err := s.setAuthorName(ctx, mitt); err != nil {
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

	for _, mitt := range mitts {
		if err := s.setLikesCount(ctx, mitt); err != nil {
			return nil, &models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
		}

		if err := s.setAuthorName(ctx, mitt); err != nil {
			return nil, &models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
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
		slog.Error("error updating mitt", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	if err := s.setLikesCount(ctx, newMitt); err != nil {
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	if err := s.setAuthorName(ctx, newMitt); err != nil {
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
		slog.Error("error deleting mitt", slog.Any("err", err))
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

func (s *Service) SwitchLike(ctx context.Context, userID uuid.UUID, mittID uuid.UUID) (bool, *models.HTTPError) {
	isAlreadyLiked, err := s.mr.IsMittLikedByUser(ctx, userID, mittID)
	if err != nil {
		slog.Error("error getting isAlreadyLiked", slog.Any("err", err))
		return false, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	if !isAlreadyLiked {
		if err := s.mr.LikeMitt(ctx, userID, mittID); err != nil {
			slog.Error("error liking mitt", slog.Any("err", err))
			return false, &models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
		}

		// Add like in metrics
		go s.mm.AddLike()

		return true, nil
	}

	if err := s.mr.DeleteMittLike(ctx, userID, mittID); err != nil {
		slog.Error("error deleting mitt like", slog.Any("err", err))
		return false, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	// Delete like in metrics
	go s.mm.DeleteLike()

	return false, nil
}

func (s *Service) Feed(ctx context.Context, limit, offset int32) ([]*models.Mitt, *models.HTTPError) {
	mitts, err := s.mr.Feed(ctx, limit, offset)
	if err != nil {
		slog.Error("failed to get feed", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	for _, mitt := range mitts {
		if err := s.setLikesCount(ctx, mitt); err != nil {
			return nil, &models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
		}

		if err := s.setAuthorName(ctx, mitt); err != nil {
			return nil, &models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
		}
	}

	return mitts, nil
}
