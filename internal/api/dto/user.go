package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Login string    `json:"login"`
	Name  string    `json:"name"`
}
