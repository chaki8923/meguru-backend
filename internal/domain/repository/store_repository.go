package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/chaki-ry/meguru-backend/internal/domain/entity"
)

type StoreRepository interface {
	Create(ctx context.Context, store *entity.Store) error
	Update(ctx context.Context, store *entity.Store) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Store, error)
}
