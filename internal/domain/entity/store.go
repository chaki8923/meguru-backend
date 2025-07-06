package entity

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Prefecture string    `json:"prefecture"`
	City       string    `json:"city"`
	Street     string    `json:"street"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
