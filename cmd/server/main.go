package main

import (
	"log"
	"os"

	infraDB "meguru-backend/internal/infrastructure/database"
	"meguru-backend/internal/infrastructure/router"
	"meguru-backend/internal/infrastructure/webpush"
	"meguru-backend/internal/interface/controller"
	"meguru-backend/internal/usecase"
	"meguru-backend/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

    
	// Database configuration
	dbConfig := database.GetConfigFromEnv()

	// Connect to database
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := infraDB.NewUserRepository(db)
	storeRepo := infraDB.NewStoreRepository(db)
	pushSubscriptionRepo := infraDB.NewPushSubscriptionRepository(db)
	flyerRepo := infraDB.NewFlyerRepository(db)

	// Initialize webpush service
	webPushService := webpush.NewWebPushService()

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo)
	healthUsecase := usecase.NewHealthUsecase()
	storeUsecase := usecase.NewStoreUsecase(storeRepo)
	pushSubscriptionUsecase := usecase.NewPushSubscriptionUsecase(pushSubscriptionRepo, webPushService)
	flyerUsecase := usecase.NewFlyerUsecase(flyerRepo)

	// Initialize controllers
	userController := controller.NewUserController(userUsecase)
	healthController := controller.NewHealthController(healthUsecase)
	storeController := controller.NewStoreController(storeUsecase)
	pushSubscriptionController := controller.NewPushSubscriptionController(pushSubscriptionUsecase)
	flyerController := controller.NewFlyerController(flyerUsecase)

	// Initialize router
	r := router.NewRouter(userController, healthController, storeController, pushSubscriptionController, flyerController)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 