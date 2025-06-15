package routes

import (
	"github.com/gin-gonic/gin"

	"meguru-backend/internal/presentation/handler"
	"meguru-backend/internal/usecase"
)

func HealthRoutes(router *gin.Engine) {
	healthUsecase := usecase.NewHealthUsecase()
	healthHandler := handler.NewHealthHandler(healthUsecase)

	router.GET("/health", healthHandler.GetHealth)
}
