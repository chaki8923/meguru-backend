package entity

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
<<<<<<< HEAD
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	Password         string     `json:"-"`
	PhoneNumber      string     `json:"phone_number"`
	Zipcode          string     `json:"zipcode"`
	Prefecture       string     `json:"prefecture"`
	City             string     `json:"city"`
	Street           string     `json:"street"`
	EmailVerifiedAt  *time.Time `json:"email_verified_at"` // メール認証済み日時（NULL許可）
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// IsEmailVerified メール認証済みかどうかをチェック
func (s *Store) IsEmailVerified() bool {
	return s.EmailVerifiedAt != nil
=======
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Prefecture string    `json:"prefecture"`
	City       string    `json:"city"`
	Street     string    `json:"street"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
>>>>>>> origin/main
}
