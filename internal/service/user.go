package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/models"
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

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.ur.GetUserByID(ctx, id)
}
