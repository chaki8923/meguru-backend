package router

import (
	"meguru-backend/internal/interface/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(userController *controller.UserController, healthController *controller.HealthController, storeController *controller.StoreController) *gin.Engine {
	r := gin.Default()

	// CORS設定
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	// ヘルスチェックエンドポイント
	r.GET("/health", healthController.GetHealth)

	// API routes
	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/register", userController.CreateUser)
		}
		stores := api.Group("/stores")
		{
			stores.POST("/", storeController.CreateStore)
			stores.PUT("/:id", storeController.UpdateStore)
		}
	}

	return r
} 