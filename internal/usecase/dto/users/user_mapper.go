package dto

import (
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/usecase/query_model"
)

func ConvertQueryModelToGetUserResponse(u *query_model.User) *GetUserResponse {
	return &GetUserResponse{
		ID:     u.ID,
		UserID: u.UserID.String(),
		Email:  u.Email,
		Name:   u.Name,
	}
}

func ConvertDomainModelToGetUserResponse(u *entity.User) *GetUserResponse {
	return &GetUserResponse{
		ID:     u.ID,
		UserID: u.UserID.String(),
		Email:  u.Email.String(),
		Name:   u.Name.String(),
	}
}
