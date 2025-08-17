
package repository

import (
    "context"

    "github.com/google/uuid"
    "meguru-backend/internal/domain/entity"
)

type TweetRepository interface {
    FindAll(ctx context.Context) ([]entity.Tweet, error)
    FindByID(ctx context.Context, id uuid.UUID) (entity.Tweet, error)
    Create(ctx context.Context, tweet entity.Tweet) (entity.Tweet, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type TweetLikeRepository interface {
    Create(ctx context.Context, tweetLike entity.TweetLike) (entity.TweetLike, error)
    Delete(ctx context.Context, tweetID, storeID uuid.UUID) error
    FindByTweetID(ctx context.Context, tweetID uuid.UUID) ([]entity.TweetLike, error)
}
