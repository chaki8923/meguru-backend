package dto

type CreateStoreRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number"`
	Zipcode     string `json:"zipcode"`
	Prefecture  string `json:"prefecture"`
	City        string `json:"city"`
	Street      string `json:"street"`
	Password    string `json:"password" binding:"required,min=6"`
}
