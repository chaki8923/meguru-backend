package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 店舗別商品情報（価格、在庫など）
type StoreProduct struct {
	ID        uuid.UUID `json:"id"`
	StoreID   uuid.UUID `json:"store_id"`
	ProductID uuid.UUID `json:"product_id"`
	Price     int       `json:"price"`          // 価格（円）
	Quantity  int       `json:"quantity"`       // 在庫数
	ImageURL  string    `json:"image_url"`      // 商品画像URL
	Status    string    `json:"status"`         // 状態（在庫あり/在庫なし）
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// リレーション
	Product *Product `json:"product,omitempty"`
} 