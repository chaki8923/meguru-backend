package repository

import (
	"context"

	"github.com/google/uuid"
	"meguru-backend/internal/domain/entity"
)

type StoreRepository interface {
	Create(ctx context.Context, store *entity.Store) error
	Update(ctx context.Context, store *entity.Store) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Store, error)
	FindAll(ctx context.Context) ([]*entity.Store, error)
}
