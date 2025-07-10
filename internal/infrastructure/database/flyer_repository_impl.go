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

func (r *FlyerRepositoryImpl) SaveFlyer(ctx context.Context, flyer *entity.Flyer, flyerData *dto.FlyerData) (*entity.Flyer, uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// 1. Save Flyer
	flyer.ID = uuid.New()
	flyer.CreatedAt = time.Now()
	flyer.UpdatedAt = flyer.CreatedAt
	_, err = tx.ExecContext(ctx, "INSERT INTO flyers (id, image_data, created_at, updated_at) VALUES ($1, $2, $3, $4)", flyer.ID, flyer.ImageData, flyer.CreatedAt, flyer.UpdatedAt)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to insert into flyers: %w", err)
	}

	// 2. Find or Create Store
	var storeID uuid.UUID
	storeInfo := flyerData.StoreInfo
	err = tx.QueryRowContext(ctx, "SELECT id FROM stores WHERE name = $1", storeInfo.Name).Scan(&storeID)
	if err == sql.ErrNoRows {
		storeID = uuid.New()
		_, err = tx.ExecContext(ctx, "INSERT INTO stores (id, name, prefecture, city, street, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())", storeID, storeInfo.Name, storeInfo.Prefecture, storeInfo.City, storeInfo.Street)
		if err != nil {
			return nil, uuid.Nil, fmt.Errorf("failed to insert store: %w", err)
		}
	} else if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to query store: %w", err)
	}

	// 3. Create Campaign
	campaignID := uuid.New()
	campaignInfo := flyerData.CampaignInfo
	startDate, _ := time.Parse("2006-01-02", campaignInfo.StartDate)
	endDate, _ := time.Parse("2006-01-02", campaignInfo.EndDate)
	_, err = tx.ExecContext(ctx, "INSERT INTO campaigns (id, flyer_id, name, start_date, end_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())", campaignID, flyer.ID, campaignInfo.Name, startDate, endDate)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to insert campaign: %w", err)
	}

	// 4. Link Campaign and Store
	_, err = tx.ExecContext(ctx, "INSERT INTO campaign_stores (campaign_id, store_id) VALUES ($1, $2)", campaignID, storeID)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to insert into campaign_stores: %w", err)
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
				return nil, uuid.Nil, fmt.Errorf("failed to insert product: %w", err)
			}
		} else if err != nil {
			return nil, uuid.Nil, fmt.Errorf("failed to query product: %w", err)
		}

		// Create FlyerItem
		_, err = tx.ExecContext(ctx, "INSERT INTO flyer_items (id, campaign_id, product_id, price_excluding_tax, price_including_tax, unit, restriction_note, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())", uuid.New(), campaignID, productID, itemData.PriceExcludingTax, itemData.PriceIncludingTax, itemData.Unit, itemData.RestrictionNote)
		if err != nil {
			return nil, uuid.Nil, fmt.Errorf("failed to insert flyer_item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Transaction commit failed: %v", err)
		return nil, uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return flyer, storeID, nil
}

func (r *FlyerRepositoryImpl) GetFlyerByStoreID(ctx context.Context, storeID string) (*entity.Flyer, *dto.FlyerData, error) {
	log.Printf("GetFlyerByStoreID called with storeID: %s", storeID)
	query := `
    SELECT
        f.id, f.image_data, f.created_at, f.updated_at,
        s.name, s.prefecture, s.city, s.street,
        c.name, c.start_date, c.end_date,
        p.name, p.category,
        fi.price_excluding_tax, fi.price_including_tax, fi.unit, fi.restriction_note
    FROM flyers f
    JOIN campaigns c ON f.id = c.flyer_id
    JOIN campaign_stores cs ON c.id = cs.campaign_id
    JOIN stores s ON cs.store_id = s.id
    LEFT JOIN flyer_items fi ON c.id = fi.campaign_id
    LEFT JOIN products p ON fi.product_id = p.id
    WHERE s.id = $1
    ORDER BY f.created_at DESC, p.name
    `

	rows, err := r.db.QueryContext(ctx, query, storeID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var flyer *entity.Flyer
	flyerData := &dto.FlyerData{}
	var items []dto.FlyerItem

	for rows.Next() {
		var flyerID uuid.UUID
		var imageData []byte
		var createdAt, updatedAt time.Time
		var storeName, storePrefecture, storeCity, storeStreet string
		var campaignName string
		var startDate, endDate time.Time
		var productName, productCategory, unit, restrictionNote sql.NullString
		var priceExcludingTax, priceIncludingTax sql.NullInt64

		if err := rows.Scan(
			&flyerID, &imageData, &createdAt, &updatedAt,
			&storeName, &storePrefecture, &storeCity, &storeStreet,
			&campaignName, &startDate, &endDate,
			&productName, &productCategory,
			&priceExcludingTax, &priceIncludingTax, &unit, &restrictionNote,
		); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if flyer == nil {
			flyer = &entity.Flyer{
				ID:        flyerID,
				ImageData: imageData,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}
			flyerData.StoreInfo.Name = storeName
			flyerData.StoreInfo.Prefecture = storePrefecture
			flyerData.StoreInfo.City = storeCity
			flyerData.StoreInfo.Street = storeStreet
			flyerData.CampaignInfo.Name = campaignName
			flyerData.CampaignInfo.StartDate = startDate.Format("2006-01-02")
			flyerData.CampaignInfo.EndDate = endDate.Format("2006-01-02")
		}

		if productName.Valid {
			items = append(items, dto.FlyerItem{
				Product: dto.Product{
					Name:     productName.String,
					Category: productCategory.String,
				},
				PriceExcludingTax: int(priceExcludingTax.Int64),
				PriceIncludingTax: int(priceIncludingTax.Int64),
				Unit:              unit.String,
				RestrictionNote:   restrictionNote.String,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if flyer == nil {
		return nil, nil, nil // No flyer found
	}

	flyerData.FlyerItemsInfo = items
	return flyer, flyerData, nil
}