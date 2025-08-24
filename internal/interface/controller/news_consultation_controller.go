package controller

import (
	"meguru-backend/internal/usecase"
	"meguru-backend/internal/usecase/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NewsConsultationController struct {
	newsConsultationUsecase *usecase.NewsConsultationUsecase
}

func NewNewsConsultationController(newsConsultationUsecase *usecase.NewsConsultationUsecase) *NewsConsultationController {
	return &NewsConsultationController{
		newsConsultationUsecase: newsConsultationUsecase,
	}
}

// ConsultNews ニュース内容からスーパーの対応策をAIに分析させる
func (c *NewsConsultationController) ConsultNews(ctx *gin.Context) {
	// リクエストボディをパース
	var request dto.NewsConsultationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// バリデーション
	if request.NewsURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "news_url is required",
		})
		return
	}

	if request.NewsTitle == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "news_title is required",
		})
		return
	}

	if request.NewsID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "news_id is required",
		})
		return
	}

	// usecase実行
	response, err := c.newsConsultationUsecase.ConsultNews(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to analyze news",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// SendNewsAnalysisEmail ニュース分析結果をメール送信
func (c *NewsConsultationController) SendNewsAnalysisEmail(ctx *gin.Context) {
	// リクエストボディをパース
	var request dto.NewsEmailRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// usecase実行
	err := c.newsConsultationUsecase.SendNewsAnalysisEmail(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send email",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Email sent successfully",
	})
}
