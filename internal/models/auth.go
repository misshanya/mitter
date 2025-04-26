package models

import "github.com/google/uuid"

type Token struct {
	Token  uuid.UUID
	UserID uuid.UUID
}

type SignIn struct {
	Login    string
	Password string
}
