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

	DeleteUser(ctx context.Context, id uuid.UUID) error

	UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdate) error
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

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) *models.HTTPError {
	// note: handle error if user not exists
	err := s.ur.DeleteUser(ctx, id)
	if err != nil {
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, user *models.UserUpdate) *models.HTTPError {
	err := s.ur.UpdateUser(ctx, id, user)
	if err != nil {
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	return nil
}
