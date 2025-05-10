package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/misshanya/mitter/internal/models"
	"github.com/misshanya/mitter/pkg/pgutil"
	"log/slog"
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

		slog.Error("error getting user", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) *models.HTTPError {
	// note: handle error if user not exists
	err := s.ur.DeleteUser(ctx, id)
	if err != nil {
		slog.Error("error deleting user", slog.Any("err", err))
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
		slog.Error("error updating user", slog.Any("err", err))
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	return nil
}

func (s *Service) FollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) *models.HTTPError {
	// Check if user (followee) exists
	_, err := s.ur.GetUserByID(ctx, followeeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &models.HTTPError{
				Code:    http.StatusNotFound,
				Message: "User not found",
			}
		}

		slog.Error("error getting user", slog.Any("err", err))
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	err = s.ur.FollowUser(ctx, followerID, followeeID)
	if err != nil {
		if pgutil.IsUniqueViolation(err) {
			return &models.HTTPError{
				Code:    http.StatusConflict,
				Message: "Already followed",
			}
		}

		slog.Error("error following user", slog.Any("err", err))
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	return nil
}

func (s *Service) UnfollowUser(ctx context.Context, followerID uuid.UUID, followeeID uuid.UUID) *models.HTTPError {
	err := s.ur.UnfollowUser(ctx, followerID, followeeID)
	if err != nil {
		slog.Error("error unfollowing user", slog.Any("err", err))
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	return nil
}

func (s *Service) GetUserFollows(ctx context.Context, followerID uuid.UUID, limit, offset int32) ([]*models.User, *models.HTTPError) {
	// Get user follows (ids)
	usersIDs, err := s.ur.GetUserFollows(ctx, followerID, limit, offset)
	if err != nil {
		slog.Error("error getting user follows", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	// Get users models from ids
	users := make([]*models.User, len(usersIDs))
	for i, id := range usersIDs {
		user, err := s.ur.GetUserByID(ctx, id)
		if err != nil {
			slog.Error("error getting user follows (getting user from db)", slog.Any("err", err))
			return nil, &models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
		}
		users[i] = user
	}

	return users, nil
}

func (s *Service) GetUserFollowers(ctx context.Context, followeeID uuid.UUID, limit, offset int32) ([]*models.User, *models.HTTPError) {
	// Get user followers (ids)
	usersIDs, err := s.ur.GetUserFollowers(ctx, followeeID, limit, offset)
	if err != nil {
		slog.Error("error getting user followers", slog.Any("err", err))
		return nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	// Get users models from ids
	users := make([]*models.User, len(usersIDs))
	for i, id := range usersIDs {
		user, err := s.ur.GetUserByID(ctx, id)
		if err != nil {
			slog.Error("error getting user followers (getting user from db)", slog.Any("err", err))
			return nil, &models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
		}
		users[i] = user
	}

	return users, nil
}
