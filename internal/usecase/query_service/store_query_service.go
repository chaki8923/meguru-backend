package query_service

import (
	"context"

	"github.com/google/uuid"

	"meguru-backend/internal/usecase/query_model"
)

type StoreQueryService interface {
	GetStoreByID(ctx context.Context, storeID uuid.UUID) (*query_model.Stores, error)
}
