package repository

import (
	"context"

	"meguru-backend/internal/domain/entity"
)

type StoreRepository interface {
	Create(ctx context.Context, store *entity.Store) error
	FindByEmail(ctx context.Context, email string) (*entity.Store, error)
}
