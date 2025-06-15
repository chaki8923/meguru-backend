package dto

type SigninUserResponse struct {
	User  *GetUserResponse `json:"user"`
	Token string           `json:"token"`
}
