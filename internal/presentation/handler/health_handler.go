package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"meguru-backend/internal/usecase"
)

type HealthHandler struct {
	healthUsecase *usecase.HealthUsecase
}

func NewHealthHandler(healthUsecase *usecase.HealthUsecase) *HealthHandler {
	return &HealthHandler{
		healthUsecase: healthUsecase,
	}
}

func (h *HealthHandler) GetHealth(c *gin.Context) {
	health := h.healthUsecase.GetHealthStatus()
	c.JSON(http.StatusOK, health)
}
