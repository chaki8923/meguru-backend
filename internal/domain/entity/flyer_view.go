package entity

import (
	"time"

	"github.com/google/uuid"
)

type FlyerView struct {
	ID       string    `json:"id"`
	FlyerID  uuid.UUID `json:"flyer_id"`
	UserID   string    `json:"user_id"`
	ViewedAt time.Time `json:"viewed_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// リレーション
	Flyer Flyer `json:"flyer,omitempty"`
	User  User  `json:"user,omitempty"`
}
