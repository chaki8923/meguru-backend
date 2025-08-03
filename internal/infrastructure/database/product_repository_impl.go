package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
)

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

// 商品マスター操作
func (r *ProductRepositoryImpl) CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt

	query := `
		INSERT INTO products (id, name, category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := r.db.ExecContext(ctx, query, product.ID, product.Name, product.Category, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (r *ProductRepositoryImpl) GetProductByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	query := `
		SELECT id, name, category, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	
	var product entity.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Category, &product.CreatedAt, &product.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	return &product, nil
}

func (r *ProductRepositoryImpl) GetProductByName(ctx context.Context, name string) (*entity.Product, error) {
	query := `
		SELECT id, name, category, created_at, updated_at
		FROM products
		WHERE name = $1
	`
	
	var product entity.Product
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&product.ID, &product.Name, &product.Category, &product.CreatedAt, &product.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product by name: %w", err)
	}

	return &product, nil
}

func (r *ProductRepositoryImpl) ListProducts(ctx context.Context) ([]*entity.Product, error) {
	query := `
		SELECT id, name, category, created_at, updated_at
		FROM products
		ORDER BY name
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	return products, nil
}

// 店舗別商品操作
func (r *ProductRepositoryImpl) CreateStoreProduct(ctx context.Context, storeProduct *entity.StoreProduct) (*entity.StoreProduct, error) {
	storeProduct.ID = uuid.New()
	storeProduct.CreatedAt = time.Now()
	storeProduct.UpdatedAt = storeProduct.CreatedAt

	query := `
		INSERT INTO store_products (id, store_id, product_id, price, quantity, image_url, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := r.db.ExecContext(ctx, query, 
		storeProduct.ID, storeProduct.StoreID, storeProduct.ProductID, 
		storeProduct.Price, storeProduct.Quantity, storeProduct.ImageURL, 
		storeProduct.Status, storeProduct.CreatedAt, storeProduct.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create store product: %w", err)
	}

	return storeProduct, nil
}

func (r *ProductRepositoryImpl) UpdateStoreProduct(ctx context.Context, storeProduct *entity.StoreProduct) (*entity.StoreProduct, error) {
	storeProduct.UpdatedAt = time.Now()

	query := `
		UPDATE store_products 
		SET price = $2, quantity = $3, image_url = $4, status = $5, updated_at = $6
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query, 
		storeProduct.ID, storeProduct.Price, storeProduct.Quantity, 
		storeProduct.ImageURL, storeProduct.Status, storeProduct.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update store product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return nil, fmt.Errorf("store product not found")
	}

	return storeProduct, nil
}

func (r *ProductRepositoryImpl) DeleteStoreProduct(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM store_products WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete store product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("store product not found")
	}

	return nil
}

func (r *ProductRepositoryImpl) GetStoreProductByID(ctx context.Context, id uuid.UUID) (*entity.StoreProduct, error) {
	query := `
		SELECT sp.id, sp.store_id, sp.product_id, sp.price, sp.quantity, sp.image_url, sp.status, sp.created_at, sp.updated_at,
		       p.id, p.name, p.category, p.created_at, p.updated_at
		FROM store_products sp
		JOIN products p ON sp.product_id = p.id
		WHERE sp.id = $1
	`
	
	var storeProduct entity.StoreProduct
	var product entity.Product
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&storeProduct.ID, &storeProduct.StoreID, &storeProduct.ProductID, 
		&storeProduct.Price, &storeProduct.Quantity, &storeProduct.ImageURL, 
		&storeProduct.Status, &storeProduct.CreatedAt, &storeProduct.UpdatedAt,
		&product.ID, &product.Name, &product.Category, &product.CreatedAt, &product.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get store product by ID: %w", err)
	}

	storeProduct.Product = &product
	return &storeProduct, nil
}

func (r *ProductRepositoryImpl) ListStoreProductsByStoreID(ctx context.Context, storeID uuid.UUID) ([]*entity.StoreProduct, error) {
	query := `
		SELECT sp.id, sp.store_id, sp.product_id, sp.price, sp.quantity, sp.image_url, sp.status, sp.created_at, sp.updated_at,
		       p.id, p.name, p.category, p.created_at, p.updated_at
		FROM store_products sp
		JOIN products p ON sp.product_id = p.id
		WHERE sp.store_id = $1
		ORDER BY p.name
	`
	
	rows, err := r.db.QueryContext(ctx, query, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list store products: %w", err)
	}
	defer rows.Close()

	var storeProducts []*entity.StoreProduct
	for rows.Next() {
		var storeProduct entity.StoreProduct
		var product entity.Product
		
		err := rows.Scan(
			&storeProduct.ID, &storeProduct.StoreID, &storeProduct.ProductID, 
			&storeProduct.Price, &storeProduct.Quantity, &storeProduct.ImageURL, 
			&storeProduct.Status, &storeProduct.CreatedAt, &storeProduct.UpdatedAt,
			&product.ID, &product.Name, &product.Category, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan store product: %w", err)
		}
		
		storeProduct.Product = &product
		storeProducts = append(storeProducts, &storeProduct)
	}

	return storeProducts, nil
}

func (r *ProductRepositoryImpl) GetStoreProductByStoreAndProduct(ctx context.Context, storeID, productID uuid.UUID) (*entity.StoreProduct, error) {
	query := `
		SELECT sp.id, sp.store_id, sp.product_id, sp.price, sp.quantity, sp.image_url, sp.status, sp.created_at, sp.updated_at,
		       p.id, p.name, p.category, p.created_at, p.updated_at
		FROM store_products sp
		JOIN products p ON sp.product_id = p.id
		WHERE sp.store_id = $1 AND sp.product_id = $2
	`
	
	var storeProduct entity.StoreProduct
	var product entity.Product
	
	err := r.db.QueryRowContext(ctx, query, storeID, productID).Scan(
		&storeProduct.ID, &storeProduct.StoreID, &storeProduct.ProductID, 
		&storeProduct.Price, &storeProduct.Quantity, &storeProduct.ImageURL, 
		&storeProduct.Status, &storeProduct.CreatedAt, &storeProduct.UpdatedAt,
		&product.ID, &product.Name, &product.Category, &product.CreatedAt, &product.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get store product: %w", err)
	}

	storeProduct.Product = &product
	return &storeProduct, nil
} 