package controller

import (
	"log"
	"net/http"

	"meguru-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type FlyerController struct {
	flyerUsecase *usecase.FlyerUsecase
}

func NewFlyerController(flyerUsecase *usecase.FlyerUsecase) *FlyerController {
	return &FlyerController{flyerUsecase: flyerUsecase}
}

func (c *FlyerController) UploadFlyer(ctx *gin.Context) {
	log.Println("UploadFlyer controller is called")
	file, err := ctx.FormFile("flyer_image")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "画像が送信されていません", "details": err.Error()})
		return
	}

	log.Println("Successfully received the flyer image")
	flyerResponse, err := c.flyerUsecase.AnalyzeAndSaveFlyer(ctx.Request.Context(), file)
	if err != nil {
		log.Printf("Error in flyer usecase: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Successfully analyzed and saved the flyer")
	ctx.JSON(http.StatusOK, gin.H{"data": flyerResponse})
}

func (c *FlyerController) GetFlyerByStoreID(ctx *gin.Context) {
	storeID := ctx.Param("store_id")

	flyerResponse, err := c.flyerUsecase.GetFlyerByStoreID(ctx.Request.Context(), storeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if flyerResponse == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No flyer found for this store"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": flyerResponse})
}
