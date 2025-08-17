package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"meguru-backend/internal/usecase"
)

type TweetController struct {
	tweetUsecase *usecase.TweetUsecase
}

func NewTweetController(tweetUsecase *usecase.TweetUsecase) *TweetController {
	return &TweetController{
		tweetUsecase: tweetUsecase,
	}
}

func (c *TweetController) ListTweets(ctx *gin.Context) {
	tweets, err := c.tweetUsecase.ListTweets(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tweets})
}

func (c *TweetController) CreateTweet(ctx *gin.Context) {
	storeID, ok := ctx.Get("storeID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req usecase.CreateTweetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tweet, err := c.tweetUsecase.CreateTweet(ctx.Request.Context(), storeID.(uuid.UUID), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": tweet})
}

func (c *TweetController) DeleteTweet(ctx *gin.Context) {
	storeID, ok := ctx.Get("storeID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tweetIDStr := ctx.Param("id")
	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet ID"})
		return
	}

	err = c.tweetUsecase.DeleteTweet(ctx.Request.Context(), storeID.(uuid.UUID), tweetID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "tweet deleted"})
}

func (c *TweetController) LikeTweet(ctx *gin.Context) {
	storeID, ok := ctx.Get("storeID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tweetIDStr := ctx.Param("id")
	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet ID"})
		return
	}

	err = c.tweetUsecase.LikeTweet(ctx.Request.Context(), storeID.(uuid.UUID), tweetID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "tweet liked"})
}

func (c *TweetController) UnlikeTweet(ctx *gin.Context) {
	storeID, ok := ctx.Get("storeID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tweetIDStr := ctx.Param("id")
	tweetID, err := uuid.Parse(tweetIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tweet ID"})
		return
	}

	err = c.tweetUsecase.UnlikeTweet(ctx.Request.Context(), storeID.(uuid.UUID), tweetID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "tweet unliked"})
}