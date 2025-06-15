package controller

import (
	"net/http"

	"meguru-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	healthUsecase *usecase.HealthUsecase
}

func NewHealthController(healthUsecase *usecase.HealthUsecase) *HealthController {
	return &HealthController{
		healthUsecase: healthUsecase,
	}
}

func (hc *HealthController) GetHealth(c *gin.Context) {
	health := hc.healthUsecase.GetHealthStatus()
	c.JSON(http.StatusOK, health)
} 