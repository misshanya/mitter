package dto

import "github.com/google/uuid"

type UserCreateRequest struct {
	Login    string `json:"login" validate:"required,min=2,max=50"`
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UserCreateResponse struct {
	ID uuid.UUID `json:"id"`
}
