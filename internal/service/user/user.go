package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/misshanya/mitter/internal/models"
	"net/http"
)

type Service struct {
	ur models.UserRepository
	um models.UserMetrics
}

func NewUserService(repo models.UserRepository, metrics models.UserMetrics) *Service {
	return &Service{
		ur: repo,
		um: metrics,
	}
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*models.User, *models.HTTPError) {
	user, err := s.ur.GetUserByID(ctx, id)
	if err != nil {
		// If user not found
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "User not found",
			}
		}
	}

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) *models.HTTPError {
	// note: handle error if user not exists
	err := s.ur.DeleteUser(ctx, id)
	if err != nil {
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	// Update metrics
	go s.um.DeleteUser()

	return nil
}

func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdate) *models.HTTPError {
	err := s.ur.UpdateUser(ctx, id, user)
	if err != nil {
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	return nil
}
