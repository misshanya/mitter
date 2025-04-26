package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/misshanya/mitter/internal/models"
	"log/slog"
	"net/http"
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

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (s *UserService) CreateUser(ctx context.Context, user models.UserCreate) (uuid.UUID, *models.HTTPError) {
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		if isUniqueViolation(err) {
			slog.Error("user already exists", slog.String("login", user.Login))
			return uuid.Nil, &models.HTTPError{
				Code:    http.StatusConflict,
				Message: "User already exists",
			}
		}
		slog.Error("error creating user", slog.Any("err", err))
		return uuid.Nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	return id, nil
}
