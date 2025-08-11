package entity

import (
	"time"

	"github.com/google/uuid"
)

type NewsView struct {
	ID        uuid.UUID `json:"id"`
	StoreID   uuid.UUID `json:"store_id"`
	NewsURL   string    `json:"news_url"`
	NewsTitle string    `json:"news_title"`
	NewsID    string    `json:"news_id"`
	ViewedAt  time.Time `json:"viewed_at"`
	CreatedAt time.Time `json:"created_at"`
}
