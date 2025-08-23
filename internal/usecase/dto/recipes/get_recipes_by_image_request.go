package dto

type GetRecipesByImageRequest struct {
	ImageBase64 string `json:"image_base64" binding:"required"`
}
