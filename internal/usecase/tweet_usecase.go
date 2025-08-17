
package usecase

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "meguru-backend/internal/domain/entity"
    "meguru-backend/internal/domain/repository"
)

type CreateTweetRequest struct {
    Content string `json:"content"`
}

type TweetResponse struct {
    ID        uuid.UUID `json:"id"`
    StoreID   uuid.UUID `json:"store_id"`
    Content   string    `json:"content"`
    Likes     int       `json:"likes"`
    CreatedAt string    `json:"created_at"`
}

type TweetUsecase struct {
    tweetRepo     repository.TweetRepository
    tweetLikeRepo repository.TweetLikeRepository
    storeRepo     repository.StoreRepository
}

func NewTweetUsecase(tweetRepo repository.TweetRepository, tweetLikeRepo repository.TweetLikeRepository, storeRepo repository.StoreRepository) *TweetUsecase {
    return &TweetUsecase{
        tweetRepo:     tweetRepo,
        tweetLikeRepo: tweetLikeRepo,
        storeRepo:     storeRepo,
    }
}

func (uc *TweetUsecase) ListTweets(ctx context.Context) ([]TweetResponse, error) {
    tweets, err := uc.tweetRepo.FindAll(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to list tweets: %w", err)
    }

    var tweetResponses []TweetResponse
    for _, tweet := range tweets {
        likes, err := uc.tweetLikeRepo.FindByTweetID(ctx, tweet.ID)
        if err != nil {
            return nil, fmt.Errorf("failed to get likes for tweet %s: %w", tweet.ID, err)
        }

        tweetResponses = append(tweetResponses, TweetResponse{
            ID:        tweet.ID,
            StoreID:   tweet.StoreID,
            Content:   tweet.Content,
            Likes:     len(likes),
            CreatedAt: tweet.CreatedAt.Format("2006-01-02 15:04:05"),
        })
    }

    return tweetResponses, nil
}

func (uc *TweetUsecase) CreateTweet(ctx context.Context, storeID uuid.UUID, req CreateTweetRequest) (TweetResponse, error) {
    if len(req.Content) > 300 {
        return TweetResponse{}, fmt.Errorf("content exceeds 300 characters")
    }

    tweet := entity.Tweet{
        StoreID: storeID,
        Content: req.Content,
    }

    createdTweet, err := uc.tweetRepo.Create(ctx, tweet)
    if err != nil {
        return TweetResponse{}, fmt.Errorf("failed to create tweet: %w", err)
    }

    return TweetResponse{
        ID:        createdTweet.ID,
        StoreID:   createdTweet.StoreID,
        Content:   createdTweet.Content,
        Likes:     0,
        CreatedAt: createdTweet.CreatedAt.Format("2006-01-02 15:04:05"),
    }, nil
}

func (uc *TweetUsecase) DeleteTweet(ctx context.Context, storeID, tweetID uuid.UUID) error {
    tweet, err := uc.tweetRepo.FindByID(ctx, tweetID)
    if err != nil {
        return fmt.Errorf("failed to find tweet: %w", err)
    }

    if tweet.StoreID != storeID {
        return fmt.Errorf("unauthorized to delete this tweet")
    }

    return uc.tweetRepo.Delete(ctx, tweetID)
}

func (uc *TweetUsecase) LikeTweet(ctx context.Context, storeID, tweetID uuid.UUID) error {
    _, err := uc.tweetRepo.FindByID(ctx, tweetID)
    if err != nil {
        return fmt.Errorf("failed to find tweet: %w", err)
    }

    like := entity.TweetLike{
        TweetID: tweetID,
        StoreID: storeID,
    }

    _, err = uc.tweetLikeRepo.Create(ctx, like)
    if err != nil {
        return fmt.Errorf("failed to like tweet: %w", err)
    }

    return nil
}

func (uc *TweetUsecase) UnlikeTweet(ctx context.Context, storeID, tweetID uuid.UUID) error {
    return uc.tweetLikeRepo.Delete(ctx, tweetID, storeID)
}
