package dto

import (
	"time"
)

type CreateStoreResponse struct {
	Store *GetStoreResponse `json:"store"`
	Token string            `json:"token"`
}

type GetStoreResponse struct {
	ID          int64     `json:"id"`
	StoreID     string    `json:"store_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Zipcode     string    `json:"zipcode"`
	Prefecture  string    `json:"prefecture"`
	City        string    `json:"city"`
	Street      string    `json:"street"`
	CreatedAt   time.Time `json:"created_at"`
}
