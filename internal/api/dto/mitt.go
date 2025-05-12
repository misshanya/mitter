package dto

import (
	"github.com/google/uuid"
	"time"
)

type MittCreateRequest struct {
	Content string `json:"content"`
}

type MittUpdateRequest struct {
	Content string `json:"content"`
}

type MittResponse struct {
	ID         uuid.UUID `json:"id"`
	Author     uuid.UUID `json:"author"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Likes      int64     `json:"likes"`
}

type MittLikeResponse struct {
	Like bool `json:"like"`
}
