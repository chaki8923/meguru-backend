
package database

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "github.com/google/uuid"
    "meguru-backend/internal/domain/entity"
    "meguru-backend/internal/domain/repository"
)

type TweetRepositoryImpl struct {
    db *sql.DB
}

func NewTweetRepository(db *sql.DB) repository.TweetRepository {
    return &TweetRepositoryImpl{db: db}
}

func (r *TweetRepositoryImpl) FindAll(ctx context.Context) ([]entity.Tweet, error) {
    query := `
        SELECT id, store_id, content, created_at, updated_at
        FROM tweets
        ORDER BY created_at DESC
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to list tweets: %w", err)
    }
    defer rows.Close()

    var tweets []entity.Tweet
    for rows.Next() {
        var tweet entity.Tweet
        err := rows.Scan(&tweet.ID, &tweet.StoreID, &tweet.Content, &tweet.CreatedAt, &tweet.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan tweet: %w", err)
        }
        tweets = append(tweets, tweet)
    }

    return tweets, nil
}

func (r *TweetRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.Tweet, error) {
    query := `
        SELECT id, store_id, content, created_at, updated_at
        FROM tweets
        WHERE id = $1
    `

    var tweet entity.Tweet
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &tweet.ID, &tweet.StoreID, &tweet.Content, &tweet.CreatedAt, &tweet.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return entity.Tweet{}, nil
    }
    if err != nil {
        return entity.Tweet{}, fmt.Errorf("failed to get tweet by ID: %w", err)
    }

    return tweet, nil
}

func (r *TweetRepositoryImpl) Create(ctx context.Context, tweet entity.Tweet) (entity.Tweet, error) {
    tweet.ID = uuid.New()
    tweet.CreatedAt = time.Now()
    tweet.UpdatedAt = tweet.CreatedAt

    query := `
        INSERT INTO tweets (id, store_id, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `

    _, err := r.db.ExecContext(ctx, query, tweet.ID, tweet.StoreID, tweet.Content, tweet.CreatedAt, tweet.UpdatedAt)
    if err != nil {
        return entity.Tweet{}, fmt.Errorf("failed to create tweet: %w", err)
    }

    return tweet, nil
}

func (r *TweetRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
    query := `DELETE FROM tweets WHERE id = $1`

    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete tweet: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("tweet not found")
    }

    return nil
}

type TweetLikeRepositoryImpl struct {
    db *sql.DB
}

func NewTweetLikeRepository(db *sql.DB) repository.TweetLikeRepository {
    return &TweetLikeRepositoryImpl{db: db}
}

func (r *TweetLikeRepositoryImpl) Create(ctx context.Context, tweetLike entity.TweetLike) (entity.TweetLike, error) {
    tweetLike.ID = uuid.New()
    tweetLike.CreatedAt = time.Now()

    query := `
        INSERT INTO tweet_likes (id, tweet_id, store_id, created_at)
        VALUES ($1, $2, $3, $4)
    `

    _, err := r.db.ExecContext(ctx, query, tweetLike.ID, tweetLike.TweetID, tweetLike.StoreID, tweetLike.CreatedAt)
    if err != nil {
        return entity.TweetLike{}, fmt.Errorf("failed to create tweet like: %w", err)
    }

    return tweetLike, nil
}

func (r *TweetLikeRepositoryImpl) Delete(ctx context.Context, tweetID, storeID uuid.UUID) error {
    query := `DELETE FROM tweet_likes WHERE tweet_id = $1 AND store_id = $2`

    result, err := r.db.ExecContext(ctx, query, tweetID, storeID)
    if err != nil {
        return fmt.Errorf("failed to delete tweet like: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("tweet like not found")
    }

    return nil
}

func (r *TweetLikeRepositoryImpl) FindByTweetID(ctx context.Context, tweetID uuid.UUID) ([]entity.TweetLike, error) {
    query := `
        SELECT id, tweet_id, store_id, created_at
        FROM tweet_likes
        WHERE tweet_id = $1
    `

    rows, err := r.db.QueryContext(ctx, query, tweetID)
    if err != nil {
        return nil, fmt.Errorf("failed to find tweet likes: %w", err)
    }
    defer rows.Close()

    var tweetLikes []entity.TweetLike
    for rows.Next() {
        var tweetLike entity.TweetLike
        err := rows.Scan(&tweetLike.ID, &tweetLike.TweetID, &tweetLike.StoreID, &tweetLike.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan tweet like: %w", err)
        }
        tweetLikes = append(tweetLikes, tweetLike)
    }

    return tweetLikes, nil
}
