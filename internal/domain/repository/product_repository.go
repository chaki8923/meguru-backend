package repository

import (
	"context"

	"github.com/google/uuid"
	"meguru-backend/internal/domain/entity"
)

type ProductRepository interface {
	// 商品マスター操作
	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
	GetProductByName(ctx context.Context, name string) (*entity.Product, error)
	ListProducts(ctx context.Context) ([]*entity.Product, error)
	
	// 店舗別商品操作
	CreateStoreProduct(ctx context.Context, storeProduct *entity.StoreProduct) (*entity.StoreProduct, error)
	UpdateStoreProduct(ctx context.Context, storeProduct *entity.StoreProduct) (*entity.StoreProduct, error)
	DeleteStoreProduct(ctx context.Context, id uuid.UUID) error
	GetStoreProductByID(ctx context.Context, id uuid.UUID) (*entity.StoreProduct, error)
	ListStoreProductsByStoreID(ctx context.Context, storeID uuid.UUID, includeExpired bool) ([]*entity.StoreProduct, error)
	GetStoreProductByStoreAndProduct(ctx context.Context, storeID, productID uuid.UUID) (*entity.StoreProduct, error)
} 