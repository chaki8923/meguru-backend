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

	// 1. Find or Create Store first (need storeID for flyer)
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

	// 2. Save Flyer with owner_store_id
	flyer.ID = uuid.New()
	flyer.OwnerStoreID = storeID
	flyer.CreatedAt = time.Now()
	flyer.UpdatedAt = flyer.CreatedAt
	
	// DisplayExpiryDateがflyerDataから提供される場合は設定
	if flyerData.DisplayExpiryDate != nil {
		flyer.DisplayExpiryDate = flyerData.DisplayExpiryDate
	}
	
	_, err = tx.ExecContext(ctx, "INSERT INTO flyers (id, image_data, owner_store_id, display_expiry_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", flyer.ID, flyer.ImageData, flyer.OwnerStoreID, flyer.DisplayExpiryDate, flyer.CreatedAt, flyer.UpdatedAt)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to insert into flyers: %w", err)
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
		return nil, uuid.Nil, fmt.Errorf("failed to insert campaign_stores: %w", err)
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

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return flyer, storeID, nil
}

// 既存の店舗IDにチラシを関連付けて保存するメソッド
func (r *FlyerRepositoryImpl) SaveFlyerForStore(ctx context.Context, flyer *entity.Flyer, flyerData *dto.FlyerData, storeID uuid.UUID) (*entity.Flyer, uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// 1. Save Flyer with owner_store_id
	flyer.ID = uuid.New()
	flyer.OwnerStoreID = storeID
	flyer.CreatedAt = time.Now()
	flyer.UpdatedAt = flyer.CreatedAt
	
	// DisplayExpiryDateがflyerDataから提供される場合は設定
	if flyerData.DisplayExpiryDate != nil {
		flyer.DisplayExpiryDate = flyerData.DisplayExpiryDate
	}
	
	_, err = tx.ExecContext(ctx, "INSERT INTO flyers (id, image_data, owner_store_id, display_expiry_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", flyer.ID, flyer.ImageData, flyer.OwnerStoreID, flyer.DisplayExpiryDate, flyer.CreatedAt, flyer.UpdatedAt)
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
}

func (r *FlyerRepositoryImpl) GetFlyerByStoreID(ctx context.Context, storeID string) (*entity.Flyer, *dto.FlyerData, error) {
	log.Printf("GetFlyerByStoreID called with storeID: %s", storeID)
	query := `
    SELECT
        f.id, f.image_data, f.display_expiry_date, f.created_at, f.updated_at,
        s.name, s.prefecture, s.city, s.street,
        c.name, c.start_date::text, c.end_date::text,
        p.name, p.category,
        fi.price_excluding_tax, fi.price_including_tax, fi.unit, fi.restriction_note
    FROM flyers f
    JOIN stores s ON f.owner_store_id = s.id
    JOIN campaigns c ON f.id = c.flyer_id
    LEFT JOIN flyer_items fi ON c.id = fi.campaign_id
    LEFT JOIN products p ON fi.product_id = p.id
    WHERE s.id = $1 
      AND (f.display_expiry_date IS NULL OR DATE(f.display_expiry_date) >= CURRENT_DATE)
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
		var displayExpiryDate sql.NullTime
		var createdAt, updatedAt time.Time
		var storeName, storePrefecture, storeCity, storeStreet string
		var campaignName string
		var startDate, endDate sql.NullTime  // time.Time から sql.NullTime に変更
		var productName, productCategory, unit, restrictionNote sql.NullString
		var priceExcludingTax, priceIncludingTax sql.NullInt64

		if err := rows.Scan(
			&flyerID, &imageData, &displayExpiryDate, &createdAt, &updatedAt,
			&storeName, &storePrefecture, &storeCity, &storeStreet,
			&campaignName, &startDate, &endDate,
			&productName, &productCategory,
			&priceExcludingTax, &priceIncludingTax, &unit, &restrictionNote,
		); err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if flyer == nil {
			var displayExpiry *time.Time
			if displayExpiryDate.Valid {
				displayExpiry = &displayExpiryDate.Time
			}
			
			flyer = &entity.Flyer{
				ID:                flyerID,
				ImageData:         imageData,
				DisplayExpiryDate: displayExpiry,
				CreatedAt:         createdAt,
				UpdatedAt:         updatedAt,
			}
			flyerData.StoreInfo.Name = storeName
			flyerData.StoreInfo.Prefecture = storePrefecture
			flyerData.StoreInfo.City = storeCity
			flyerData.StoreInfo.Street = storeStreet
			flyerData.CampaignInfo.Name = campaignName
			flyerData.DisplayExpiryDate = displayExpiry
			
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
}

// GetAllFlyersByStoreID retrieves all flyers for a specific store ID
func (r *FlyerRepositoryImpl) GetAllFlyersByStoreID(ctx context.Context, storeID string) ([]*entity.Flyer, []*dto.FlyerData, error) {
	log.Printf("FlyerRepository: Getting all flyers for storeID: %s", storeID)

	// Use the same query structure as GetFlyerByStoreID but remove LIMIT 1
	query := `
		SELECT 
			f.id, f.image_data, f.display_expiry_date, f.created_at, f.updated_at,
			s.name, s.prefecture, s.city, s.street,
			c.name, c.start_date::text, c.end_date::text,
			p.name, p.category,
			fi.price_excluding_tax, fi.price_including_tax, fi.unit, fi.restriction_note
		FROM flyers f
		JOIN stores s ON f.owner_store_id = s.id
		JOIN campaigns c ON f.id = c.flyer_id
		LEFT JOIN flyer_items fi ON c.id = fi.campaign_id
		LEFT JOIN products p ON fi.product_id = p.id
		WHERE f.owner_store_id = $1
		  AND (f.display_expiry_date IS NULL OR DATE(f.display_expiry_date) >= CURRENT_DATE)
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
		var displayExpiryDate sql.NullTime
		var createdAt, updatedAt time.Time
		var campaignName sql.NullString
		var startDate, endDate sql.NullTime
		var productName, productCategory, unit, restrictionNote sql.NullString
		var priceExcludingTax, priceIncludingTax sql.NullInt64

		err := rows.Scan(
			&flyerID, &imageData, &displayExpiryDate, &createdAt, &updatedAt,
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

			var displayExpiry *time.Time
			if displayExpiryDate.Valid {
				displayExpiry = &displayExpiryDate.Time
			}

			flyerMap[flyerID] = &entity.Flyer{
				ID:                parsedID,
				ImageData:         []byte(imageData),
				DisplayExpiryDate: displayExpiry,
				CreatedAt:         createdAt,
				UpdatedAt:         updatedAt,
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
				FlyerItemsInfo:    []dto.FlyerItem{},
				DisplayExpiryDate: displayExpiry,
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
}

// GetNearbyFlyers 近隣店舗のチラシを取得
func (r *FlyerRepositoryImpl) GetNearbyFlyers(ctx context.Context, city string, limit int) ([]*entity.Flyer, []*dto.FlyerData, error) {
	log.Printf("FlyerRepository: Getting nearby flyers for city: %s, limit: %d", city, limit)
	
	// デバッグ用: 全チラシ数と指定都市の店舗数を確認
	var totalFlyers, cityStores int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM flyers").Scan(&totalFlyers)
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stores WHERE city = $1", city).Scan(&cityStores)
	log.Printf("Debug: Total flyers in DB: %d, Stores in city '%s': %d", totalFlyers, city, cityStores)

	query := `
		SELECT
			f.id, f.image_data, f.display_expiry_date, f.created_at, f.updated_at,
			owner_store.id, owner_store.name, owner_store.prefecture, owner_store.city, owner_store.street,
			COALESCE(c.name, ''), COALESCE(c.start_date::text, ''), COALESCE(c.end_date::text, ''),
			COALESCE(p.name, ''), COALESCE(p.category, ''), COALESCE(fi.price_excluding_tax, 0), COALESCE(fi.price_including_tax, 0), COALESCE(fi.unit, ''), COALESCE(fi.restriction_note, '')
		FROM flyers f
		JOIN stores owner_store ON f.owner_store_id = owner_store.id
		JOIN campaigns c ON f.id = c.flyer_id
		LEFT JOIN flyer_items fi ON c.id = fi.campaign_id
		LEFT JOIN products p ON fi.product_id = p.id
		WHERE owner_store.city = $1
		  AND (f.display_expiry_date IS NULL OR DATE(f.display_expiry_date) >= CURRENT_DATE)
		ORDER BY f.created_at DESC, p.name
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, city, limit*10) // より多く取得して後でフィルタリング
	if err != nil {
		log.Printf("Error executing nearby flyers query: %v", err)
		return nil, nil, fmt.Errorf("error querying nearby flyers: %w", err)
	}
	defer rows.Close()
	
	log.Printf("Query executed successfully for city: %s", city)

	flyerMap := make(map[uuid.UUID]*entity.Flyer)
	flyerDataMap := make(map[uuid.UUID]*dto.FlyerData)
	flyerOrder := []uuid.UUID{}

	for rows.Next() {
		var flyerID uuid.UUID
		var imageData []byte
		var displayExpiryDate sql.NullTime
		var createdAt, updatedAt time.Time

		var storeID uuid.UUID
		var storeName, prefecture, city, street string

		var campaignName, startDate, endDate string
		var productName, category string
		var priceExcludingTax, priceIncludingTax int
		var unit, restrictionNote string

		if err := rows.Scan(
			&flyerID, &imageData, &displayExpiryDate, &createdAt, &updatedAt,
			&storeID, &storeName, &prefecture, &city, &street,
			&campaignName, &startDate, &endDate,
			&productName, &category, &priceExcludingTax, &priceIncludingTax, &unit, &restrictionNote,
		); err != nil {
			return nil, nil, fmt.Errorf("error scanning row: %w", err)
		}

		log.Printf("Debug row: flyerID=%s, storeID=%s, storeName=%s, campaignName=%s", flyerID, storeID, storeName, campaignName)

		if _, exists := flyerMap[flyerID]; !exists {
			var displayExpiry *time.Time
			if displayExpiryDate.Valid {
				displayExpiry = &displayExpiryDate.Time
			}

			flyerMap[flyerID] = &entity.Flyer{
				ID:                flyerID,
				ImageData:         imageData,
				DisplayExpiryDate: displayExpiry,
				CreatedAt:         createdAt,
				UpdatedAt:         updatedAt,
			}

			flyerDataMap[flyerID] = &dto.FlyerData{
				StoreInfo: dto.Store{
					ID:         storeID.String(),
					Name:       storeName,
					Prefecture: prefecture,
					City:       city,
					Street:     street,
				},
				CampaignInfo: dto.Campaign{
					Name:      campaignName,
					StartDate: startDate,
					EndDate:   endDate,
				},
				FlyerItemsInfo:    []dto.FlyerItem{},
				DisplayExpiryDate: displayExpiry,
			}

			flyerOrder = append(flyerOrder, flyerID)
		}

		// 商品情報を追加（productNameが空でない場合のみ）
		if productName != "" {
			flyerItem := dto.FlyerItem{
				Product: dto.Product{
					Name:     productName,
					Category: category,
				},
				PriceExcludingTax: priceExcludingTax,
				PriceIncludingTax: priceIncludingTax,
				Unit:              unit,
				RestrictionNote:   restrictionNote,
			}
			flyerDataMap[flyerID].FlyerItemsInfo = append(flyerDataMap[flyerID].FlyerItemsInfo, flyerItem)
		}

		// 制限数に達したらbreak
		if len(flyerMap) >= limit {
			break
		}
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if len(flyerMap) == 0 {
		return []*entity.Flyer{}, []*dto.FlyerData{}, nil
	}

	// Convert maps to slices in order (制限数まで)
	actualLimit := limit
	if len(flyerOrder) < limit {
		actualLimit = len(flyerOrder)
	}
	
	flyers := make([]*entity.Flyer, actualLimit)
	flyerDataList := make([]*dto.FlyerData, actualLimit)
	
	for i := 0; i < actualLimit; i++ {
		flyerID := flyerOrder[i]
		flyers[i] = flyerMap[flyerID]
		flyerDataList[i] = flyerDataMap[flyerID]
	}

	log.Printf("FlyerRepository: Successfully retrieved %d nearby flyers for city: %s", len(flyers), city)
	return flyers, flyerDataList, nil
}