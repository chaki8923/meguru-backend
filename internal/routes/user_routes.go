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

	userRouterGroup := router.Group("/api/v1").Group("/users")

	userRouterGroup.POST("/signup", userHandler.CreateUser)
	userRouterGroup.POST("/signin", userHandler.Signin)
	userRouterGroup.GET("/:user_id", middleware.ValidateJWTMiddleware(), userHandler.GetUserByID)

	return router
}
