package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"meguru-backend/internal/domain/domain_service"
	"meguru-backend/internal/infrastructure/query_service"
	infraDB "meguru-backend/internal/infrastructure/repository"
	"meguru-backend/internal/middleware"
	"meguru-backend/internal/presentation/handler"
	"meguru-backend/internal/usecase"
)

func UserRoutes(db *sql.DB, router *gin.Engine) *gin.Engine {
	userRepo := infraDB.NewUserRepository(db)
	userQueryService := query_service.NewUserQueryService(db)
	userDomainService := domain_service.NewUserDomainService(userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, userDomainService, userQueryService)
	userHandler := handler.NewUserHandler(userUsecase)

	router.Group("/api/v1").Group("/users").POST("/signup", userHandler.CreateUser)
	router.Group("/api/v1").Group("/users").POST("/signin", userHandler.Signin)
	router.Group("/api/v1").Group("/users").GET("/:user_id", middleware.ValidateJWTMiddleware(), userHandler.GetUserByID)

	return router
}
