package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"meguru-backend/internal/domain/domain_service"
	"meguru-backend/internal/infrastructure/repository"
	"meguru-backend/internal/presentation/handler"
	"meguru-backend/internal/usecase"
)

func StoreRoutes(db *sql.DB, router *gin.Engine) *gin.Engine {
	storeRepo := repository.NewStoreRepository(db)
	storeDomainService := domain_service.NewStoreService(storeRepo)
	storeUsecase := usecase.NewStoreUsecase(storeRepo, storeDomainService)
	storeHandler := handler.NewStoreHandler(storeUsecase)

	router.Group("/api/v1").Group("/stores").POST("/signup", storeHandler.CreateStore)
	router.Group("/api/v1").Group("/stores").GET("", storeHandler.GetStores)

	return router
}
