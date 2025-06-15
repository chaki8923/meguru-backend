package dto

import (
	"meguru-backend/internal/usecase/query_model"
)

func MapToGetUserResponse(u *query_model.User) *GetUserResponse {
	return &GetUserResponse{
		ID:        u.ID,
		UserID:    u.UserID.String(),
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}
