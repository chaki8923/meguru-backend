package entity

import (
	"time"

	"github.com/google/uuid"
)

type Flyer struct {
	ID                uuid.UUID  `json:"id"`
	ImageData         []byte     `json:"image_data"`
	OwnerStoreID      uuid.UUID  `json:"owner_store_id"`
	DisplayExpiryDate *time.Time `json:"display_expiry_date,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
