package models

import (
	"context"
	"github.com/google/uuid"
)

type MittRepository interface {
	CreateMitt(ctx context.Context, userID uuid.UUID, mitt *MittCreate) (*Mitt, error)

	GetMitt(ctx context.Context, id uuid.UUID) (*Mitt, error)
	GetAllUserMitts(ctx context.Context, userID uuid.UUID) ([]*Mitt, error)

	UpdateMitt(ctx context.Context, mittID uuid.UUID, mitt *MittUpdate) (*Mitt, error)

	DeleteMitt(ctx context.Context, mittID uuid.UUID) error
}
