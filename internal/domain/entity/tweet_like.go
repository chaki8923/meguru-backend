
package entity

import (
    "time"

    "github.com/google/uuid"
)

type TweetLike struct {
    ID        uuid.UUID `json:"id"`
    TweetID   uuid.UUID `json:"tweet_id"`
    StoreID   uuid.UUID `json:"store_id"`
    CreatedAt time.Time `json:"created_at"`
}
