package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/misshanya/mitter/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
)

type authRepository interface {
	SaveToken(ctx context.Context, token *models.Token) error
	GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error)
}

type AuthService struct {
	ur userRepository
	ar authRepository
}

func NewAuthService(ur userRepository, ar authRepository) *AuthService {
	return &AuthService{ur: ur, ar: ar}
}

func (s *AuthService) SignIn(ctx context.Context, creds models.SignIn) (string, *models.HTTPError) {
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

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func (s *AuthService) SignUp(ctx context.Context, user *models.UserCreate) (uuid.UUID, *models.HTTPError) {
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

func (s *AuthService) ChangePassword(ctx context.Context, id uuid.UUID, changePassword *models.ChangePassword) *models.HTTPError {
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
