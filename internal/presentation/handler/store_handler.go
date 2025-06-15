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

	resp, err := uc.storeUsecase.CreateStore(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (uc *StoreHandler) SigninStore(c *gin.Context) {
	var req dto.SigninStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := uc.storeUsecase.SigninStore(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (uc *StoreHandler) GetStoreByID(c *gin.Context) {
	storeID := c.Param("store_id")

	store, err := uc.storeUsecase.GetStoreByID(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if store == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "store not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": store})
}
