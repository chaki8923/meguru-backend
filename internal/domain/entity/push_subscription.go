package entity

import (
	"time"
)

type PushSubscription struct {
	ID        int       `json:"id"`
	Endpoint  string    `json:"endpoint"`
	P256dh    string    `json:"p256dh"`
	Auth      string    `json:"auth"`
	CreatedAt time.Time `json:"created_at"`
}