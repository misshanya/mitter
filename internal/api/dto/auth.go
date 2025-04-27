package dto

import "github.com/google/uuid"

type SignUpRequest struct {
	Login    string `json:"login" validate:"required,min=2,max=50"`
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type SignUpResponse struct {
	ID uuid.UUID `json:"id"`
}

type SignInRequest struct {
	Login    string `json:"login" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8,max=100"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}
