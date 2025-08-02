package database

import (
	"context"
	"database/sql"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
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
		INSERT INTO stores (id, name, prefecture, city, street, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		store.ID, store.Name, store.Prefecture, store.City, store.Street, store.CreatedAt, store.UpdatedAt)

	return err
}

func (r *storeRepository) Update(ctx context.Context, store *entity.Store) error {
	query := `
		UPDATE stores
		SET name = $2, prefecture = $3, city = $4, street = $5, updated_at = $6
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		store.ID, store.Name, store.Prefecture, store.City, store.Street, store.UpdatedAt)

	return err
}

func (r *storeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Store, error) {
	query := `
		SELECT id, name, prefecture, city, street, created_at, updated_at
		FROM stores
		WHERE id = $1`

	store := &entity.Store{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&store.ID, &store.Name, &store.Prefecture, &store.City, &store.Street, &store.CreatedAt, &store.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (r *storeRepository) FindAll(ctx context.Context) ([]*entity.Store, error) {
	query := `
		SELECT id, name, prefecture, city, street, created_at, updated_at
		FROM stores`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := make([]*entity.Store, 0)
	for rows.Next() {
		store := &entity.Store{}
		if err := rows.Scan(
			&store.ID, &store.Name, &store.Prefecture, &store.City, &store.Street, &store.CreatedAt, &store.UpdatedAt);
		 err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stores, nil
}
