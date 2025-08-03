package repository

import (
	"context"
	"time"

	"meguru-backend/internal/domain/entity"
	"github.com/google/uuid"
)

// StoreEmailVerificationTokenRepository 店舗メール認証トークンのリポジトリインターフェース
type StoreEmailVerificationTokenRepository interface {
	// Create トークンを作成
	Create(ctx context.Context, storeID uuid.UUID, token string, expiresAt time.Time) error
	
	// FindByToken トークンで検索
	FindByToken(ctx context.Context, token string) (*entity.StoreEmailVerificationToken, error)
	
	// Delete トークンを削除
	Delete(ctx context.Context, id uuid.UUID) error
	
	// DeleteByStoreID 指定した店舗のトークンをすべて削除
	DeleteByStoreID(ctx context.Context, storeID uuid.UUID) error
	
	// DeleteExpired 期限切れトークンを削除
	DeleteExpired(ctx context.Context) error
} 