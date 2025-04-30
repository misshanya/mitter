package models

import (
	"github.com/google/uuid"
	"time"
)

type MittCreate struct {
	Content string
}

type Mitt struct {
	ID        uuid.UUID
	Author    uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Likes     int64
}

type MittUpdate struct {
	Content string
}
