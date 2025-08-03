package controller

import (
	"net/http"
	"strings"

	"meguru-backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"errors"
)

type StoreController struct {
	storeUsecase *usecase.StoreUsecase
}

func NewStoreController(storeUsecase *usecase.StoreUsecase) *StoreController {
	return &StoreController{
		storeUsecase: storeUsecase,
	}
}

func (sc *StoreController) CreateStore(c *gin.Context) {
	var req usecase.CreateStoreRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store, err := sc.storeUsecase.CreateStore(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": store})
}

func (sc *StoreController) UpdateStore(c *gin.Context) {
	var req usecase.UpdateStoreRequest

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store, err := sc.storeUsecase.UpdateStore(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": store})
}

func (sc *StoreController) GetAllStores(c *gin.Context) {
	stores, err := sc.storeUsecase.GetAllStores(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stores})
}

// 店舗登録用の新しいハンドラーメソッド
func (sc *StoreController) RegisterShop(c *gin.Context) {
	var req usecase.ShopRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := sc.storeUsecase.RegisterShop(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// 店舗ログイン用のハンドラーメソッド
func (sc *StoreController) SignIn(c *gin.Context) {
	var req usecase.StoreSignInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := sc.storeUsecase.SignIn(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Authorization ヘッダーからトークンを取得するヘルパー関数
func (sc *StoreController) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	// "Bearer " プレフィックスを削除
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), nil
	}

	return "", errors.New("invalid authorization header format")
}

// 現在ログイン中の店舗プロフィール取得
func (sc *StoreController) GetProfile(c *gin.Context) {
	token, err := sc.getTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := sc.storeUsecase.GetProfile(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile})
}

// 現在ログイン中の店舗プロフィール更新
func (sc *StoreController) UpdateProfile(c *gin.Context) {
	token, err := sc.getTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req usecase.StoreUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := sc.storeUsecase.UpdateProfile(c.Request.Context(), token, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile})
}
