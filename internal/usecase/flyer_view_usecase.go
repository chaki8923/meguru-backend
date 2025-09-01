package usecase

import (
	"context"
	"time"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"github.com/google/uuid"
)

type FlyerViewUsecase struct {
	flyerViewRepo repository.FlyerViewRepository
}

func NewFlyerViewUsecase(flyerViewRepo repository.FlyerViewRepository) *FlyerViewUsecase {
	return &FlyerViewUsecase{
		flyerViewRepo: flyerViewRepo,
	}
}

// フライヤービュー記録用のリクエスト構造体
type RecordFlyerViewRequest struct {
	FlyerID string `json:"flyer_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

// フライヤービュー記録用のレスポンス構造体
type RecordFlyerViewResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// フライヤービュー数取得用のレスポンス構造体
type FlyerViewCountResponse struct {
	FlyerID   string `json:"flyer_id"`
	ViewCount int64  `json:"view_count"`
}

// フライヤービューリスト取得用のレスポンス構造体
type FlyerViewListResponse struct {
	FlyerID string               `json:"flyer_id"`
	Views   []FlyerViewWithUser  `json:"views"`
	Total   int64                `json:"total"`
}

type FlyerViewWithUser struct {
	ID       string    `json:"id"`
	UserName string    `json:"user_name"`
	UserEmail string   `json:"user_email"`
	ViewedAt time.Time `json:"viewed_at"`
}

// フライヤーのビューを記録する
func (u *FlyerViewUsecase) RecordFlyerView(ctx context.Context, req *RecordFlyerViewRequest) (*RecordFlyerViewResponse, error) {
	// UUIDパース
	flyerID, err := uuid.Parse(req.FlyerID)
	if err != nil {
		return nil, err
	}

	// フライヤービューエンティティを作成
	flyerView := &entity.FlyerView{
		ID:        uuid.New().String(),
		FlyerID:   flyerID,
		UserID:    req.UserID,
		ViewedAt:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// ビューを記録（重複の場合は無視）
	err = u.flyerViewRepo.CreateOrIgnore(ctx, flyerView)
	if err != nil {
		return nil, err
	}

	return &RecordFlyerViewResponse{
		Success: true,
		Message: "フライヤーのビューが記録されました",
	}, nil
}

// フライヤーのビュー数を取得する
func (u *FlyerViewUsecase) GetFlyerViewCount(ctx context.Context, flyerID string) (*FlyerViewCountResponse, error) {
	// UUIDパース
	parsedFlyerID, err := uuid.Parse(flyerID)
	if err != nil {
		return nil, err
	}

	// ビュー数を取得
	count, err := u.flyerViewRepo.GetViewCountByFlyerID(ctx, parsedFlyerID)
	if err != nil {
		return nil, err
	}

	return &FlyerViewCountResponse{
		FlyerID:   flyerID,
		ViewCount: count,
	}, nil
}

// フライヤーのビューリストを取得する（店舗側での詳細確認用）
func (u *FlyerViewUsecase) GetFlyerViewList(ctx context.Context, flyerID string, limit int, offset int) (*FlyerViewListResponse, error) {
	// UUIDパース
	parsedFlyerID, err := uuid.Parse(flyerID)
	if err != nil {
		return nil, err
	}

	// ビューリストを取得
	flyerViews, err := u.flyerViewRepo.GetViewsByFlyerID(ctx, parsedFlyerID, limit, offset)
	if err != nil {
		return nil, err
	}

	// ビュー数も取得
	total, err := u.flyerViewRepo.GetViewCountByFlyerID(ctx, parsedFlyerID)
	if err != nil {
		return nil, err
	}

	// レスポンス用に変換
	var views []FlyerViewWithUser
	for _, fv := range flyerViews {
		views = append(views, FlyerViewWithUser{
			ID:        fv.ID,
			UserName:  fv.User.Name,
			UserEmail: fv.User.Email,
			ViewedAt:  fv.ViewedAt,
		})
	}

	return &FlyerViewListResponse{
		FlyerID: flyerID,
		Views:   views,
		Total:   total,
	}, nil
}
