package dto

type SigninStoreResponse struct {
	Token string            `json:"token"`
	Store *GetStoreResponse `json:"store"`
}
