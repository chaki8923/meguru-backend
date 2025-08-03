package controller

import (
	"log"
	"net/http"
	"strings"
	"errors"

	"meguru-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type FlyerController struct {
	flyerUsecase *usecase.FlyerUsecase
}

func NewFlyerController(flyerUsecase *usecase.FlyerUsecase) *FlyerController {
	return &FlyerController{flyerUsecase: flyerUsecase}
}

// Authorization ヘッダーからトークンを取得するヘルパー関数
func (fc *FlyerController) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	log.Printf("Received Authorization header: '%s'", authHeader) // デバッグログ
	
	if authHeader == "" {
		log.Println("Authorization header is empty") // デバッグログ
		return "", errors.New("authorization header required")
	}

	// "Bearer " プレフィックスを削除
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Extracted token: '%s'", token) // デバッグログ
		return token, nil
	}

	log.Println("Invalid authorization header format") // デバッグログ
	return "", errors.New("invalid authorization header format")
}

func (c *FlyerController) UploadFlyer(ctx *gin.Context) {
	log.Println("UploadFlyer controller is called")

	// 認証チェック（UpdateProfileと同じパターン）
	token, err := c.getTokenFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := ctx.FormFile("flyer_image")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "画像が送信されていません", "details": err.Error()})
		return
	}

	log.Println("Successfully received the flyer image")
	
	// ログイン中の店舗情報を更新
	flyerResponse, err := c.flyerUsecase.AnalyzeAndUpdateStoreFromFlyer(ctx.Request.Context(), file, token)
	if err != nil {
		log.Printf("Error in flyer usecase: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Successfully analyzed flyer and updated store information")
	ctx.JSON(http.StatusOK, gin.H{"data": flyerResponse})
}

func (c *FlyerController) GetFlyerByStoreID(ctx *gin.Context) {
	storeID := ctx.Param("store_id")
	log.Printf("GetFlyerByStoreID called with storeID: %s", storeID)

	// panic回復
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in GetFlyerByStoreID: %v", r)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}()

	if c.flyerUsecase == nil {
		log.Printf("FlyerUsecase is nil")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Usecase not initialized"})
		return
	}

	log.Printf("Calling FlyerUsecase.GetFlyerByStoreID")
	flyerResponse, err := c.flyerUsecase.GetFlyerByStoreID(ctx.Request.Context(), storeID)
	if err != nil {
		log.Printf("Error in GetFlyerByStoreID usecase: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if flyerResponse == nil {
		log.Printf("No flyer found for store ID: %s", storeID)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No flyer found for this store"})
		return
	}

	log.Printf("Successfully retrieved flyer for store ID: %s", storeID)
	ctx.JSON(http.StatusOK, gin.H{"data": flyerResponse})
}

// GetAllFlyersByStoreID retrieves all flyers for a specific store ID
func (c *FlyerController) GetAllFlyersByStoreID(ctx *gin.Context) {
	storeID := ctx.Param("store_id")
	log.Printf("GetAllFlyersByStoreID called with storeID: %s", storeID)

	// panic回復
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in GetAllFlyersByStoreID: %v", r)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}()

	if c.flyerUsecase == nil {
		log.Printf("FlyerUsecase is nil")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Usecase not initialized"})
		return
	}

	log.Printf("Calling FlyerUsecase.GetAllFlyersByStoreID")
	flyerResponses, err := c.flyerUsecase.GetAllFlyersByStoreID(ctx.Request.Context(), storeID)
	if err != nil {
		log.Printf("Error in GetAllFlyersByStoreID usecase: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(flyerResponses) == 0 {
		log.Printf("No flyers found for store ID: %s", storeID)
		ctx.JSON(http.StatusOK, gin.H{"data": []interface{}{}}) // Return empty array instead of 404
		return
	}

	log.Printf("Successfully retrieved %d flyers for store ID: %s", len(flyerResponses), storeID)
	ctx.JSON(http.StatusOK, gin.H{"data": flyerResponses})
}
