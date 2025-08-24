package repository

import (
	"context"

	"meguru-backend/internal/domain/entity"
	"github.com/google/uuid"
)

type FlyerViewRepository interface {
	// フライヤーのビューを記録する（既存の場合は何もしない）
	CreateOrIgnore(ctx context.Context, flyerView *entity.FlyerView) error
	
	// 特定のフライヤーのビュー数を取得する
	GetViewCountByFlyerID(ctx context.Context, flyerID uuid.UUID) (int64, error)
	
	// 特定のフライヤーのビューリストを取得する（店舗側での詳細確認用）
	GetViewsByFlyerID(ctx context.Context, flyerID uuid.UUID, limit int, offset int) ([]*entity.FlyerView, error)
	
	// ユーザーが特定のフライヤーを既に見たかどうかを確認
	HasUserViewedFlyer(ctx context.Context, flyerID uuid.UUID, userID string) (bool, error)
}
