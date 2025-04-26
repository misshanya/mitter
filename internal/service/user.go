package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/misshanya/mitter/internal/models"
	"net/http"
)

type userRepository interface {
	CreateUser(ctx context.Context, user *models.UserCreate) (uuid.UUID, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type UserService struct {
	ur userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{ur: repo}
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, *models.HTTPError) {
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
