package controller

import (
	"meguru-backend/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type NewsViewController struct {
	newsViewUsecase *usecase.NewsViewUsecase
}

func NewNewsViewController(newsViewUsecase *usecase.NewsViewUsecase) *NewsViewController {
	return &NewsViewController{
		newsViewUsecase: newsViewUsecase,
	}
}

// RecordNewsView ニュース閲覧を記録
func (c *NewsViewController) RecordNewsView(ctx *gin.Context) {
	// トークンを取得
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
		return
	}

	// リクエストボディをパース
	var request usecase.RecordNewsViewRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// ニュース閲覧を記録
	err := c.newsViewUsecase.RecordNewsView(ctx.Request.Context(), token, &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record news view", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "News view recorded successfully"})
}

// GetNewsViewCount 単一ニュースの閲覧数を取得
func (c *NewsViewController) GetNewsViewCount(ctx *gin.Context) {
	newsID := ctx.Param("news_id")
	if newsID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "News ID is required"})
		return
	}

	response, err := c.newsViewUsecase.GetNewsViewCount(ctx.Request.Context(), newsID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news view count", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetNewsViewCounts 複数ニュースの閲覧数を取得
func (c *NewsViewController) GetNewsViewCounts(ctx *gin.Context) {
	var request struct {
		NewsIDs []string `json:"news_ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	counts, err := c.newsViewUsecase.GetNewsViewCounts(ctx.Request.Context(), request.NewsIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news view counts", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"view_counts": counts})
}
