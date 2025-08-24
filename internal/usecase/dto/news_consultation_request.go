package dto

// NewsConsultationRequest ニュース相談リクエスト
type NewsConsultationRequest struct {
	NewsURL   string `json:"news_url" binding:"required"`
	NewsTitle string `json:"news_title" binding:"required"`
	NewsID    string `json:"news_id" binding:"required"`
}
