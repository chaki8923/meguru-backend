package entity

import (
	"time"
	"github.com/google/uuid"
)

// StoreEmailVerificationToken 店舗メール認証トークン
type StoreEmailVerificationToken struct {
	ID        uuid.UUID `json:"id"`
	StoreID   uuid.UUID `json:"store_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// IsExpired トークンが期限切れかどうかをチェック
func (t *StoreEmailVerificationToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
} 