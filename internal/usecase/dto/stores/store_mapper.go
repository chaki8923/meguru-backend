package dto

import (
	"meguru-backend/internal/usecase/query_model"
)

func MapToGetStoresResponse(s *query_model.Stores) *GetStoreResponse {
	return &GetStoreResponse{
		ID:          s.ID,
		StoreID:     s.StoreID.String(),
		Name:        s.Name,
		Email:       s.Email,
		PhoneNumber: s.PhoneNumber,
		Zipcode:     s.Zipcode,
		Prefecture:  s.Prefecture,
		City:        s.City,
		Street:      s.Street,
		CreatedAt:   s.CreatedAt,
	}
}
