package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"

	"github.com/google/uuid"
)

type NewsViewUsecase struct {
	newsViewRepository repository.NewsViewRepository
	storeRepository    repository.StoreRepository
}

func NewNewsViewUsecase(newsViewRepository repository.NewsViewRepository, storeRepository repository.StoreRepository) *NewsViewUsecase {
	return &NewsViewUsecase{
		newsViewRepository: newsViewRepository,
		storeRepository:    storeRepository,
	}
}

type RecordNewsViewRequest struct {
	NewsURL   string `json:"news_url" binding:"required"`
	NewsTitle string `json:"news_title" binding:"required"`
	NewsID    string `json:"news_id" binding:"required"`
}

type NewsViewCountResponse struct {
	NewsID    string `json:"news_id"`
	ViewCount int    `json:"view_count"`
}

// RecordNewsView ニュース閲覧を記録
func (u *NewsViewUsecase) RecordNewsView(ctx context.Context, token string, request *RecordNewsViewRequest) error {
	// トークンから店舗IDを取得
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		log.Printf("Error getting store ID from token: %v", err)
		return fmt.Errorf("failed to get store ID: %w", err)
	}

	// 店舗の存在確認
	_, err = u.storeRepository.FindByID(ctx, storeID)
	if err != nil {
		log.Printf("Store not found: %s", storeID)
		return fmt.Errorf("store not found")
	}

	// ニュース閲覧記録を作成
	newsView := &entity.NewsView{
		StoreID:   storeID,
		NewsURL:   request.NewsURL,
		NewsTitle: request.NewsTitle,
		NewsID:    request.NewsID,
	}

	_, err = u.newsViewRepository.CreateNewsView(ctx, newsView)
	if err != nil {
		// 30分以内の重複閲覧の場合は正常終了
		if err.Error() == "news already viewed within 30 minutes" {
			log.Printf("Duplicate view within 30 minutes, ignoring: store=%s, news=%s", storeID, request.NewsID)
			return nil
		}
		log.Printf("Error creating news view: %v", err)
		return fmt.Errorf("failed to record news view: %w", err)
	}

	log.Printf("News view recorded successfully: store=%s, news=%s", storeID, request.NewsID)
	return nil
}

// GetNewsViewCount 単一ニュースの閲覧数を取得
func (u *NewsViewUsecase) GetNewsViewCount(ctx context.Context, newsID string) (*NewsViewCountResponse, error) {
	count, err := u.newsViewRepository.GetNewsViewCount(ctx, newsID)
	if err != nil {
		log.Printf("Error getting news view count for %s: %v", newsID, err)
		return nil, fmt.Errorf("failed to get news view count: %w", err)
	}

	return &NewsViewCountResponse{
		NewsID:    newsID,
		ViewCount: count,
	}, nil
}

// GetNewsViewCounts 複数ニュースの閲覧数を取得
func (u *NewsViewUsecase) GetNewsViewCounts(ctx context.Context, newsIDs []string) (map[string]int, error) {
	counts, err := u.newsViewRepository.GetNewsViewCountsByNewsIDs(ctx, newsIDs)
	if err != nil {
		log.Printf("Error getting news view counts: %v", err)
		return nil, fmt.Errorf("failed to get news view counts: %w", err)
	}

	return counts, nil
}

// トークンからストアIDを取得するヘルパー関数
func (u *NewsViewUsecase) getStoreIDFromToken(token string) (uuid.UUID, error) {
	// productユースケースと同じロジックを使用
	var uuidStr string
	if strings.HasPrefix(token, "auth_token_") {
		uuidStr = strings.TrimPrefix(token, "auth_token_")
	} else if strings.HasPrefix(token, "temp_token_") {
		uuidStr = strings.TrimPrefix(token, "temp_token_")
	} else {
		return uuid.Nil, fmt.Errorf("invalid token format")
	}

	storeID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	return storeID, nil
}
