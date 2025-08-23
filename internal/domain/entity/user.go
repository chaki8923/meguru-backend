package entity

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Zipcode      string    `json:"zipcode"`
	Prefecture   string    `json:"prefecture"`
	City         string    `json:"city"`
	Street       string    `json:"street"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
} 