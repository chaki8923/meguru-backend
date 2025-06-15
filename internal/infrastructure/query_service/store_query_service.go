package query_service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"meguru-backend/internal/usecase/query_model"
	"meguru-backend/internal/usecase/query_service"
)

type storeQueryService struct {
	db *sql.DB
}

func NewStoreQueryService(db *sql.DB) query_service.StoreQueryService {
	return &storeQueryService{db: db}
}

func (s *storeQueryService) GetStoreByID(ctx context.Context, storeID uuid.UUID) (*query_model.Stores, error) {
	query := `
		SELECT id, store_id, name, email, phone_number, zipcode, prefecture, city, street, created_at
		FROM stores
		WHERE store_id = $1
	`
	row := s.db.QueryRowContext(ctx, query, storeID.String())

	store := &query_model.Stores{}
	var storeIDStr string

	err := row.Scan(
		&store.ID,
		&storeIDStr,
		&store.Name,
		&store.Email,
		&store.PhoneNumber,
		&store.Zipcode,
		&store.Prefecture,
		&store.City,
		&store.Street,
		&store.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	storeIDParsed, err := uuid.Parse(storeIDStr)
	if err != nil {
		return nil, err
	}
	store.StoreID = storeIDParsed

	return store, nil
}
