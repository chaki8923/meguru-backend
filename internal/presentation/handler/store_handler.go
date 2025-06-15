package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"meguru-backend/internal/usecase"
	dto "meguru-backend/internal/usecase/dto/stores"
)

type StoreHandler struct {
	storeUsecase *usecase.StoreUsecase
}

func NewStoreHandler(storeUsecase *usecase.StoreUsecase) *StoreHandler {
	return &StoreHandler{
		storeUsecase: storeUsecase,
	}
}

func (uc *StoreHandler) CreateStore(c *gin.Context) {
	var req dto.CreateStoreRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store, token, err := uc.storeUsecase.CreateStore(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": store, "token": token})
}

func (uc *StoreHandler) GetStores(c *gin.Context) {
	stores, err := uc.storeUsecase.GetStores(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stores})
}
