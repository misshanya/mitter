package models

type UserCreate struct {
	Login          string
	Name           string
	Password       string
	HashedPassword string
}
