package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/misshanya/mitter/internal/models"
)

type userRepository interface {
	CreateUser(ctx context.Context, user models.UserCreate) (uuid.UUID, error)
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user models.UserCreate) (uuid.UUID, error) {
	return s.repo.CreateUser(ctx, user)
}
