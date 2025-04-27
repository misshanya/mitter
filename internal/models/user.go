package models

import "github.com/google/uuid"

type UserCreate struct {
	Login          string
	Name           string
	Password       string
	HashedPassword string
}

type User struct {
	ID             uuid.UUID
	Login          string
	Name           string
	HashedPassword string
}

type UserUpdate struct {
	Name *string
}
