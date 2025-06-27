package database

import (
	"context"
	"database/sql"

	"github.com/chaki-ry/meguru-backend/internal/domain/entity"
	"github.com/chaki-ry/meguru-backend/internal/domain/repository"
	"github.com/google/uuid"
)

type storeRepository struct {
	db *sql.DB
}

func NewStoreRepository(db *sql.DB) repository.StoreRepository {
	return &storeRepository{
		db: db,
	}
}

func (r *storeRepository) Create(ctx context.Context, store *entity.Store) error {
	query := `
		INSERT INTO stores (id, name, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		store.ID, store.Name, store.Address, store.CreatedAt, store.UpdatedAt)

	return err
}

func (r *storeRepository) Update(ctx context.Context, store *entity.Store) error {
	query := `
		UPDATE stores
		SET name = $2, address = $3, updated_at = $4
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		store.ID, store.Name, store.Address, store.UpdatedAt)

	return err
}

func (r *storeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Store, error) {
	query := `
		SELECT id, name, address, created_at, updated_at
		FROM stores
		WHERE id = $1`

	store := &entity.Store{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&store.ID, &store.Name, &store.Address, &store.CreatedAt, &store.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return store, nil
}
