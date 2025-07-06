package handler

import (
	"meguru-backend/internal/presentation/responses"
	"meguru-backend/internal/usecase"
	dto "meguru-backend/internal/usecase/dto/stores"

	"github.com/gin-gonic/gin"
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
		responses.HTTP400(c, gin.H{"error": err.Error()})
		return
	}

	resp, err := uc.storeUsecase.CreateStore(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "store with this email already exists" {
			responses.HTTP409(c, gin.H{"error": err.Error()})
			return
		}

		responses.HTTP500(c, gin.H{"error": "internal server error"})
		return
	}

	responses.HTTP201(c, gin.H{"data": resp})
}

func (uc *StoreHandler) SigninStore(c *gin.Context) {
	var req dto.SigninStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.HTTP400(c, gin.H{"error": err.Error()})
		return
	}

	resp, err := uc.storeUsecase.SigninStore(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			responses.HTTP401(c, gin.H{"error": err.Error()})
			return
		}

		responses.HTTP500(c, gin.H{"error": "internal server error"})
		return
	}

	responses.HTTP200(c, gin.H{"data": resp})
}

func (uc *StoreHandler) GetStoreByID(c *gin.Context) {
	storeID := c.Param("store_id")

	store, err := uc.storeUsecase.GetStoreByID(c.Request.Context(), storeID)
	if err != nil {
		responses.HTTP500(c, gin.H{"error": err.Error()})
		return
	}

	if store == nil {
		responses.HTTP404(c, gin.H{"error": "store not found"})
		return
	}

	responses.HTTP200(c, gin.H{"data": store})
}
