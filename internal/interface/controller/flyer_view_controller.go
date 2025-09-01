package controller

import (
	"net/http"
	"strconv"

	"meguru-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

type FlyerViewController struct {
	flyerViewUsecase *usecase.FlyerViewUsecase
}

func NewFlyerViewController(flyerViewUsecase *usecase.FlyerViewUsecase) *FlyerViewController {
	return &FlyerViewController{
		flyerViewUsecase: flyerViewUsecase,
	}
}

// フライヤーのビューを記録する
func (fvc *FlyerViewController) RecordFlyerView(c *gin.Context) {
	var req usecase.RecordFlyerViewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := fvc.flyerViewUsecase.RecordFlyerView(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// フライヤーのビュー数を取得する
func (fvc *FlyerViewController) GetFlyerViewCount(c *gin.Context) {
	flyerID := c.Param("flyer_id")
	if flyerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "flyer_id is required"})
		return
	}

	response, err := fvc.flyerViewUsecase.GetFlyerViewCount(c.Request.Context(), flyerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// フライヤーのビューリストを取得する（店舗側での詳細確認用）
func (fvc *FlyerViewController) GetFlyerViewList(c *gin.Context) {
	flyerID := c.Param("flyer_id")
	if flyerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "flyer_id is required"})
		return
	}

	// クエリパラメータから limit と offset を取得
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	response, err := fvc.flyerViewUsecase.GetFlyerViewList(c.Request.Context(), flyerID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}
