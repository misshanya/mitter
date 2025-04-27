package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Login string    `json:"login"`
	Name  string    `json:"name"`
}

type UserUpdateRequest struct {
	Name *string `json:"name" validate:"omitempty,min=2,max=50"`
}
