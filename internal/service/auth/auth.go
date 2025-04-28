package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/misshanya/mitter/internal/models"
	"github.com/misshanya/mitter/internal/service/utils"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
)

type Service struct {
	ur models.UserRepository
	ar models.AuthRepository
}

func NewAuthService(ur models.UserRepository, ar models.AuthRepository) *Service {
	return &Service{ur: ur, ar: ar}
}

func (s *Service) SignIn(ctx context.Context, creds models.SignIn) (string, *models.HTTPError) {
	user, err := s.ur.GetUserByLogin(ctx, creds.Login)
	if err != nil {
		// If user not found
		if errors.Is(err, pgx.ErrNoRows) {
			return "", &models.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid login or password",
			}
		}

		return "", &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(creds.Password)); err != nil {
		return "", &models.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid login or password",
		}
	}

	// Generate access token
	token := uuid.New()

	// Save token to Redis
	if err := s.ar.SaveToken(ctx, &models.Token{
		Token:  token,
		UserID: user.ID,
	}); err != nil {
		slog.Error("error saving token to redis", slog.Any("err", err))
		return "", &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	return token.String(), nil
}

func (s *Service) SignUp(ctx context.Context, user *models.UserCreate) (uuid.UUID, *models.HTTPError) {
	// Hash password and store it in user.HashedPassword
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("error hashing password", slog.Any("err", err))
		return uuid.Nil, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	user.HashedPassword = string(hashedPassword)

	id, err := s.ur.CreateUser(ctx, user)
	if err != nil {
		if utils.IsUniqueViolation(err) {
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

func (s *Service) ChangePassword(ctx context.Context, id uuid.UUID, changePassword *models.ChangePassword) *models.HTTPError {
	// Compare old passwords
	currentPwdHash, err := s.ur.GetCurrentPasswordHash(ctx, id)
	if err != nil {
		// If user not found
		if errors.Is(err, pgx.ErrNoRows) {
			return &models.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "User does not exist",
			}
		}

		slog.Error("error getting current password hash", slog.Any("err", err))
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentPwdHash), []byte(changePassword.OldPassword)); err != nil {
		return &models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Old password doesn't match",
		}
	}

	// Hash new password
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(changePassword.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	// Write new hash in db
	if err := s.ur.ChangePassword(ctx, id, string(newPasswordHash)); err != nil {
		return &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	return nil
}
