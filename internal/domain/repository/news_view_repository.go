package repository

import (
	"context"
	"meguru-backend/internal/domain/entity"
	"github.com/google/uuid"
)

type NewsViewRepository interface {
	// CreateNewsView 新しいニュース閲覧記録を作成
	CreateNewsView(ctx context.Context, newsView *entity.NewsView) (*entity.NewsView, error)
	
	// CheckRecentView 指定した店舗が指定したニュースを最近閲覧したかチェック（30分以内）
	CheckRecentView(ctx context.Context, storeID uuid.UUID, newsID string) (bool, error)
	
	// GetNewsViewCount 指定したニュースの総閲覧数を取得
	GetNewsViewCount(ctx context.Context, newsID string) (int, error)
	
	// GetNewsViewCountsByNewsIDs 複数のニュースIDの閲覧数をまとめて取得
	GetNewsViewCountsByNewsIDs(ctx context.Context, newsIDs []string) (map[string]int, error)
}
