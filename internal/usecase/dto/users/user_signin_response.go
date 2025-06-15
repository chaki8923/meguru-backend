package dto

type SigninResponse struct {
	Token string           `json:"token"`
	User  *GetUserResponse `json:"user"`
}
