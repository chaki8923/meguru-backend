package dto

import (
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/usecase/query_model"
)

func ConvertQueryModelToGetUserResponse(u *query_model.User) *GetUserResponse {
	return &GetUserResponse{
		ID:        u.ID,
		UserID:    u.UserID.String(),
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}

func ConvertDomainModelToGetUserResponse(user *entity.User) *GetUserResponse {
	return &GetUserResponse{
		ID:        user.ID,
		UserID:    user.UserID.String(),
		Email:     user.Email.String(),
		Name:      user.Name.String(),
		CreatedAt: user.CreatedAt,
	}
}
