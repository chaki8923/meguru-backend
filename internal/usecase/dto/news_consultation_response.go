package dto

// NewsConsultationResponse ニュース相談レスポンス
type NewsConsultationResponse struct {
	NewsID          string `json:"news_id"`
	NewsTitle       string `json:"news_title"`
	NewsURL         string `json:"news_url"`
	AnalysisResult  string `json:"analysis_result"`
	Recommendations string `json:"recommendations"`
	CreatedAt       string `json:"created_at"`
}
