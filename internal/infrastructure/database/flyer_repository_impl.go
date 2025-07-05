package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/dto"

	"github.com/google/uuid"
)

type FlyerRepositoryImpl struct {
	db *sql.DB
}

func NewFlyerRepository(db *sql.DB) repository.FlyerRepository {
	return &FlyerRepositoryImpl{db: db}
}

func (r *FlyerRepositoryImpl) SaveFlyer(ctx context.Context, flyer *entity.Flyer, flyerData *dto.FlyerData) (*entity.Flyer, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// 1. Save Flyer
	flyer.ID = uuid.New()
	flyer.CreatedAt = time.Now()
	flyer.UpdatedAt = flyer.CreatedAt
	_, err = tx.ExecContext(ctx, "INSERT INTO flyers (id, image_data, created_at, updated_at) VALUES ($1, $2, $3, $4)", flyer.ID, flyer.ImageData, flyer.CreatedAt, flyer.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into flyers: %w", err)
	}

	// 2. Find or Create Store
	var storeID uuid.UUID
	storeInfo := flyerData.StoreInfo
	err = tx.QueryRowContext(ctx, "SELECT id FROM stores WHERE name = $1", storeInfo.Name).Scan(&storeID)
	if err == sql.ErrNoRows {
		storeID = uuid.New()
		_, err = tx.ExecContext(ctx, "INSERT INTO stores (id, name, address, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())", storeID, storeInfo.Name, storeInfo.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to insert store: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to query store: %w", err)
	}

	// 3. Create Campaign
	campaignID := uuid.New()
	campaignInfo := flyerData.CampaignInfo
	startDate, _ := time.Parse("2006-01-02", campaignInfo.StartDate)
	endDate, _ := time.Parse("2006-01-02", campaignInfo.EndDate)
	_, err = tx.ExecContext(ctx, "INSERT INTO campaigns (id, flyer_id, name, start_date, end_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())", campaignID, flyer.ID, campaignInfo.Name, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to insert campaign: %w", err)
	}

	// 4. Link Campaign and Store
	_, err = tx.ExecContext(ctx, "INSERT INTO campaign_stores (campaign_id, store_id) VALUES ($1, $2)", campaignID, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert into campaign_stores: %w", err)
	}

	// 5. Find or Create Products and create FlyerItems
	for _, itemData := range flyerData.FlyerItemsInfo {
		var productID uuid.UUID
		productInfo := itemData.Product
		err = tx.QueryRowContext(ctx, "SELECT id FROM products WHERE name = $1", productInfo.Name).Scan(&productID)
		if err == sql.ErrNoRows {
			productID = uuid.New()
			_, err = tx.ExecContext(ctx, "INSERT INTO products (id, name, category, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())", productID, productInfo.Name, productInfo.Category)
			if err != nil {
				return nil, fmt.Errorf("failed to insert product: %w", err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to query product: %w", err)
		}

		// Create FlyerItem
		_, err = tx.ExecContext(ctx, "INSERT INTO flyer_items (id, campaign_id, product_id, price_excluding_tax, price_including_tax, unit, restriction_note, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())", uuid.New(), campaignID, productID, itemData.PriceExcludingTax, itemData.PriceIncludingTax, itemData.Unit, itemData.RestrictionNote)
		if err != nil {
			return nil, fmt.Errorf("failed to insert flyer_item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Transaction commit failed: %v", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return flyer, nil
}

