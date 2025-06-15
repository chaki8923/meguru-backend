package dto

import (
	"time"
)

type CreateUserResponse struct {
	User  *GetUserResponse `json:"user"`
	Token string           `json:"token"`
}

type GetUserResponse struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
