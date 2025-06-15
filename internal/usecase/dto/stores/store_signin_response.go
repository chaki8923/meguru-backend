package dto

type SigninStoreResponse struct {
	Store *GetStoreResponse `json:"store"`
	Token string            `json:"token"`
}
