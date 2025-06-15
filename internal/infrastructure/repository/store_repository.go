package repository

import (
	"context"
	"database/sql"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	store_vo "meguru-backend/internal/domain/value_object/store"
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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO stores (
			store_id,
			name,
			email,
			password_hash,
			phone_number,
			zipcode,
			prefecture,
			city,
			street,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		store.StoreID.Value().String(),
		store.Name,
		store.Email.String(),
		store.PasswordHash,
		store.PhoneNumber,
		store.Zipcode,
		store.Prefecture,
		store.City,
		store.Street,
		store.CreatedAt,
		store.UpdatedAt,
	).Scan(&store.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *storeRepository) FindByEmail(ctx context.Context, email string) (*entity.Store, error) {
	query := `
		SELECT 
			id, store_id, name, email, password_hash,
			phone_number, zipcode, prefecture, city, street, created_at, updated_at
		FROM stores
		WHERE email = $1`

	store := &entity.Store{}

	var storeIDStr string
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&store.ID,
		&storeIDStr,
		&store.Name,
		&store.Email,
		&store.PasswordHash,
		&store.PhoneNumber,
		&store.Zipcode,
		&store.Prefecture,
		&store.City,
		&store.Street,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	storeIDVO, err := store_vo.NewUuid(storeIDStr)
	if err != nil {
		return nil, err
	}
	store.StoreID = *storeIDVO

	return store, nil
}
