package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type NewsViewRepositoryImpl struct {
	db *sql.DB
}

func NewNewsViewRepository(db *sql.DB) repository.NewsViewRepository {
	return &NewsViewRepositoryImpl{db: db}
}

// CreateNewsView 新しいニュース閲覧記録を作成
func (r *NewsViewRepositoryImpl) CreateNewsView(ctx context.Context, newsView *entity.NewsView) (*entity.NewsView, error) {
	// 30分以内の重複チェック
	recentView, err := r.CheckRecentView(ctx, newsView.StoreID, newsView.NewsID)
	if err != nil {
		log.Printf("Error checking recent view: %v", err)
		return nil, fmt.Errorf("failed to check recent view: %w", err)
	}
	
	if recentView {
		log.Printf("Store %s already viewed news %s within 30 minutes, skipping", newsView.StoreID, newsView.NewsID)
		return nil, fmt.Errorf("news already viewed within 30 minutes")
	}

	// IDと時刻を設定
	newsView.ID = uuid.New()
	newsView.ViewedAt = time.Now()
	newsView.CreatedAt = time.Now()

	query := `
		INSERT INTO news_views (id, store_id, news_url, news_title, news_id, viewed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (store_id, news_id) DO UPDATE SET
			viewed_at = EXCLUDED.viewed_at
		RETURNING id, store_id, news_url, news_title, news_id, viewed_at, created_at
	`

	row := r.db.QueryRowContext(ctx, query,
		newsView.ID, newsView.StoreID, newsView.NewsURL, newsView.NewsTitle, newsView.NewsID,
		newsView.ViewedAt, newsView.CreatedAt)

	var result entity.NewsView
	err = row.Scan(&result.ID, &result.StoreID, &result.NewsURL, &result.NewsTitle, &result.NewsID,
		&result.ViewedAt, &result.CreatedAt)
	if err != nil {
		log.Printf("Error inserting news view: %v", err)
		return nil, fmt.Errorf("failed to insert news view: %w", err)
	}

	log.Printf("News view recorded: store=%s, news=%s", newsView.StoreID, newsView.NewsID)
	return &result, nil
}

// CheckRecentView 指定した店舗が指定したニュースを最近閲覧したかチェック（30分以内）
func (r *NewsViewRepositoryImpl) CheckRecentView(ctx context.Context, storeID uuid.UUID, newsID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM news_views 
			WHERE store_id = $1 AND news_id = $2 
			AND viewed_at > NOW() - INTERVAL '30 minutes'
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, storeID, newsID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking recent view: %v", err)
		return false, fmt.Errorf("failed to check recent view: %w", err)
	}

	return exists, nil
}

// GetNewsViewCount 指定したニュースの総閲覧数を取得
func (r *NewsViewRepositoryImpl) GetNewsViewCount(ctx context.Context, newsID string) (int, error) {
	query := `SELECT COUNT(*) FROM news_views WHERE news_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, newsID).Scan(&count)
	if err != nil {
		log.Printf("Error getting news view count for %s: %v", newsID, err)
		return 0, fmt.Errorf("failed to get news view count: %w", err)
	}

	return count, nil
}

// GetNewsViewCountsByNewsIDs 複数のニュースIDの閲覧数をまとめて取得
func (r *NewsViewRepositoryImpl) GetNewsViewCountsByNewsIDs(ctx context.Context, newsIDs []string) (map[string]int, error) {
	if len(newsIDs) == 0 {
		return make(map[string]int), nil
	}

	// PostgreSQLの配列型として渡す
	query := `
		SELECT news_id, COUNT(*) as view_count
		FROM news_views
		WHERE news_id = ANY($1)
		GROUP BY news_id
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(newsIDs))
	if err != nil {
		log.Printf("Error getting news view counts: %v", err)
		return nil, fmt.Errorf("failed to get news view counts: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int)
	for _, newsID := range newsIDs {
		result[newsID] = 0 // 初期値として0を設定
	}

	for rows.Next() {
		var newsID string
		var count int
		if err := rows.Scan(&newsID, &count); err != nil {
			log.Printf("Error scanning news view count row: %v", err)
			continue
		}
		result[newsID] = count
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating news view count rows: %v", err)
		return nil, fmt.Errorf("failed to iterate news view count rows: %w", err)
	}

	return result, nil
}
