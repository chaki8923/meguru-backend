
package entity

import (
    "time"

    "github.com/google/uuid"
)

type Tweet struct {
    ID        uuid.UUID `json:"id"`
    StoreID   uuid.UUID `json:"store_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
