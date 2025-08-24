package dto

// NewsEmailRequest ニュース分析結果メール送信リクエスト
type NewsEmailRequest struct {
	Email           string `json:"email" binding:"required,email"`
	NewsID          string `json:"news_id" binding:"required"`
	NewsTitle       string `json:"news_title" binding:"required"`
	NewsURL         string `json:"news_url" binding:"required"`
	AnalysisResult  string `json:"analysis_result" binding:"required"`
	Recommendations string `json:"recommendations" binding:"required"`
}
