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

	// 3. Save Campaign
	campaignInfo := flyerData.CampaignInfo
	campaignID := uuid.New()
	
	// 日付文字列をパース（空文字列の場合はnilを使用）
	var startDate, endDate interface{}
	if campaignInfo.StartDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", campaignInfo.StartDate); err == nil {
			startDate = parsedStartDate
		} else {
			startDate = nil
		}
	} else {
		startDate = nil
	}
	
	if campaignInfo.EndDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", campaignInfo.EndDate); err == nil {
			endDate = parsedEndDate
		} else {
			endDate = nil
		}
	} else {
		endDate = nil
	}
	
	_, err = tx.ExecContext(ctx, "INSERT INTO campaigns (id, flyer_id, name, start_date, end_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())", campaignID, flyer.ID, campaignInfo.Name, startDate, endDate)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to insert campaign: %w", err)
	}

	// 4. Associate Campaign with Store
	_, err = tx.ExecContext(ctx, "INSERT INTO campaign_stores (campaign_id, store_id) VALUES ($1, $2)", campaignID, storeID)
	if err != nil {
<<<<<<< HEAD
		return nil, uuid.Nil, fmt.Errorf("failed to insert campaign_stores: %w", err)
=======
		return nil, uuid.Nil, fmt.Errorf("failed to insert into campaign_stores: %w", err)
>>>>>>> origin/main
	}

	// 5. Save Products and Flyer Items
	for _, item := range flyerData.FlyerItemsInfo {
		// Find or Create Product
		var productID uuid.UUID
		err = tx.QueryRowContext(ctx, "SELECT id FROM products WHERE name = $1", item.Product.Name).Scan(&productID)
		if err == sql.ErrNoRows {
			productID = uuid.New()
			_, err = tx.ExecContext(ctx, "INSERT INTO products (id, name, category, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())", productID, item.Product.Name, item.Product.Category)
			if err != nil {
				return nil, uuid.Nil, fmt.Errorf("failed to insert product: %w", err)
			}
		} else if err != nil {
			return nil, uuid.Nil, fmt.Errorf("failed to query product: %w", err)
		}

		// Save Flyer Item
		flyerItemID := uuid.New()
		_, err = tx.ExecContext(ctx, "INSERT INTO flyer_items (id, campaign_id, product_id, price_excluding_tax, price_including_tax, unit, restriction_note, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())", flyerItemID, campaignID, productID, item.PriceExcludingTax, item.PriceIncludingTax, item.Unit, item.RestrictionNote)
		if err != nil {
			return nil, uuid.Nil, fmt.Errorf("failed to insert flyer_item: %w", err)
		}
	}

<<<<<<< HEAD
	// Commit transaction
	if err = tx.Commit(); err != nil {
=======
	if err := tx.Commit(); err != nil {
		log.Printf("Transaction commit failed: %v", err)
>>>>>>> origin/main
		return nil, uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return flyer, storeID, nil
<<<<<<< HEAD
}

// 既存の店舗IDにチラシを関連付けて保存するメソッド
func (r *FlyerRepositoryImpl) SaveFlyerForStore(ctx context.Context, flyer *entity.Flyer, flyerData *dto.FlyerData, storeID uuid.UUID) (*entity.Flyer, uuid.UUID, error) {
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

	// 2. Save Campaign (既存の店舗IDを使用)
	campaignInfo := flyerData.CampaignInfo
	campaignID := uuid.New()
	
	// 日付文字列をパース（空文字列の場合はnilを使用）
	var startDate, endDate interface{}
	if campaignInfo.StartDate != "" {
		if parsedStartDate, err := time.Parse("2006-01-02", campaignInfo.StartDate); err == nil {
			startDate = parsedStartDate
		} else {
			startDate = nil
		}
	} else {
		startDate = nil
	}
	
	if campaignInfo.EndDate != "" {
		if parsedEndDate, err := time.Parse("2006-01-02", campaignInfo.EndDate); err == nil {
			endDate = parsedEndDate
		} else {
			endDate = nil
		}
	} else {
		endDate = nil
	}
	
	_, err = tx.ExecContext(ctx, "INSERT INTO campaigns (id, flyer_id, name, start_date, end_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())", campaignID, flyer.ID, campaignInfo.Name, startDate, endDate)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to insert campaign: %w", err)
	}

	// 3. Associate Campaign with existing Store
	_, err = tx.ExecContext(ctx, "INSERT INTO campaign_stores (campaign_id, store_id) VALUES ($1, $2)", campaignID, storeID)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to insert campaign_stores: %w", err)
	}

	// 4. Save Products and Flyer Items
	for _, item := range flyerData.FlyerItemsInfo {
		// Find or Create Product
		var productID uuid.UUID
		err = tx.QueryRowContext(ctx, "SELECT id FROM products WHERE name = $1", item.Product.Name).Scan(&productID)
		if err == sql.ErrNoRows {
			productID = uuid.New()
			_, err = tx.ExecContext(ctx, "INSERT INTO products (id, name, category, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())", productID, item.Product.Name, item.Product.Category)
			if err != nil {
				return nil, uuid.Nil, fmt.Errorf("failed to insert product: %w", err)
			}
		} else if err != nil {
			return nil, uuid.Nil, fmt.Errorf("failed to query product: %w", err)
		}

		// Save Flyer Item
		flyerItemID := uuid.New()
		_, err = tx.ExecContext(ctx, "INSERT INTO flyer_items (id, campaign_id, product_id, price_excluding_tax, price_including_tax, unit, restriction_note, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())", flyerItemID, campaignID, productID, item.PriceExcludingTax, item.PriceIncludingTax, item.Unit, item.RestrictionNote)
		if err != nil {
			return nil, uuid.Nil, fmt.Errorf("failed to insert flyer_item: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return flyer, storeID, nil
=======
>>>>>>> origin/main
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
		var startDate, endDate sql.NullTime  // time.Time から sql.NullTime に変更
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
			
			// sql.NullTime を適切に処理
			if startDate.Valid {
				flyerData.CampaignInfo.StartDate = startDate.Time.Format("2006-01-02")
			} else {
				flyerData.CampaignInfo.StartDate = ""
			}
			
			if endDate.Valid {
				flyerData.CampaignInfo.EndDate = endDate.Time.Format("2006-01-02")
			} else {
				flyerData.CampaignInfo.EndDate = ""
			}
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
<<<<<<< HEAD
}

// GetAllFlyersByStoreID retrieves all flyers for a specific store ID
func (r *FlyerRepositoryImpl) GetAllFlyersByStoreID(ctx context.Context, storeID string) ([]*entity.Flyer, []*dto.FlyerData, error) {
	log.Printf("FlyerRepository: Getting all flyers for storeID: %s", storeID)

	// Use the same query structure as GetFlyerByStoreID but remove LIMIT 1
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
		WHERE cs.store_id = $1
		ORDER BY f.created_at DESC, p.name
	`

	rows, err := r.db.QueryContext(ctx, query, storeID)
	if err != nil {
		log.Printf("FlyerRepository: Failed to execute query for storeID %s: %v", storeID, err)
		return nil, nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	flyerMap := make(map[string]*entity.Flyer)
	flyerDataMap := make(map[string]*dto.FlyerData)
	flyerOrder := []string{} // To maintain order

	for rows.Next() {
		log.Printf("FlyerRepository: Processing row for storeID %s", storeID)

		var flyerID, imageData, storeName, prefecture, city, street string
		var createdAt, updatedAt time.Time
		var campaignName sql.NullString
		var startDate, endDate sql.NullTime
		var productName, productCategory, unit, restrictionNote sql.NullString
		var priceExcludingTax, priceIncludingTax sql.NullInt64

		err := rows.Scan(
			&flyerID, &imageData, &createdAt, &updatedAt,
			&storeName, &prefecture, &city, &street,
			&campaignName, &startDate, &endDate,
			&productName, &productCategory,
			&priceExcludingTax, &priceIncludingTax,
			&unit, &restrictionNote,
		)
		if err != nil {
			log.Printf("FlyerRepository: Failed to scan row for storeID %s: %v", storeID, err)
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Initialize flyer if not exists
		if _, exists := flyerMap[flyerID]; !exists {
			parsedID, err := uuid.Parse(flyerID)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse flyer ID: %w", err)
			}

			flyerMap[flyerID] = &entity.Flyer{
				ID:        parsedID,
				ImageData: []byte(imageData),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}

			flyerDataMap[flyerID] = &dto.FlyerData{
				StoreInfo: dto.Store{
					Name:       storeName,
					Prefecture: prefecture,
					City:       city,
					Street:     street,
				},
				CampaignInfo: dto.Campaign{
					Name: "",
					StartDate: "",
					EndDate: "",
				},
				FlyerItemsInfo: []dto.FlyerItem{},
			}

			flyerOrder = append(flyerOrder, flyerID)
		}

		flyerData := flyerDataMap[flyerID]

		// Set campaign info (only once per flyer)
		if campaignName.Valid && flyerData.CampaignInfo.Name == "" {
			flyerData.CampaignInfo.Name = campaignName.String
			if startDate.Valid {
				flyerData.CampaignInfo.StartDate = startDate.Time.Format("2006-01-02")
			} else {
				flyerData.CampaignInfo.StartDate = ""
			}
			if endDate.Valid {
				flyerData.CampaignInfo.EndDate = endDate.Time.Format("2006-01-02")
			} else {
				flyerData.CampaignInfo.EndDate = ""
			}
		}

		// Add product if exists
		if productName.Valid {
			flyerData.FlyerItemsInfo = append(flyerData.FlyerItemsInfo, dto.FlyerItem{
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

	if len(flyerMap) == 0 {
		return []*entity.Flyer{}, []*dto.FlyerData{}, nil // No flyers found
	}

	// Convert maps to slices in order
	flyers := make([]*entity.Flyer, len(flyerOrder))
	flyerDataList := make([]*dto.FlyerData, len(flyerOrder))
	
	for i, flyerID := range flyerOrder {
		flyers[i] = flyerMap[flyerID]
		flyerDataList[i] = flyerDataMap[flyerID]
	}

	log.Printf("FlyerRepository: Successfully retrieved %d flyers for storeID: %s", len(flyers), storeID)
	return flyers, flyerDataList, nil
=======
>>>>>>> origin/main
}