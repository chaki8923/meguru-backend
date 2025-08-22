package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	infraDB "meguru-backend/internal/infrastructure/database"
	"meguru-backend/internal/infrastructure/email"
	"meguru-backend/internal/infrastructure/router"
	"meguru-backend/internal/infrastructure/storage"
	"meguru-backend/internal/infrastructure/webpush"
	"meguru-backend/internal/interface/controller"
	"meguru-backend/internal/usecase"
	"meguru-backend/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

    
	// Database configuration
	dbConfig := database.GetConfigFromEnv()

	var db *sql.DB
	var err error
	maxRetries := 15
	retryInterval := 60

	for i := 0; i < maxRetries; i++ {
		log.Printf("データベースに接続を試みています... (試行 %d/%d)\n", i+1, maxRetries)
		db, err = database.NewPostgresDB(dbConfig)
		if err == nil {
			log.Println("データベースへの接続に成功しました！")
			break
		}

		log.Printf("接続に失敗しました: %v\n", err)
		if i < maxRetries-1 {
			wait := time.Second * time.Duration(retryInterval)
			log.Printf("%v 後に再試行します。\n", wait)
			time.Sleep(wait)
		}
	}

	if err != nil {
		log.Fatalf("%d回試行しましたが、データベースに接続できませんでした: %v", maxRetries, err)
	}

	defer db.Close()

	userRepo := infraDB.NewUserRepository(db)
	passwordResetTokenRepo := infraDB.NewPasswordResetTokenRepository(db)
	storeRepo := infraDB.NewStoreRepository(db)
	storeTokenRepo := infraDB.NewStoreEmailVerificationTokenRepository(db)
	pushSubscriptionRepo := infraDB.NewPushSubscriptionRepository(db)
	flyerRepo := infraDB.NewFlyerRepository(db)
	productRepo := infraDB.NewProductRepository(db)
	newsViewRepo := infraDB.NewNewsViewRepository(db)
	tweetRepo := infraDB.NewTweetRepository(db)
	tweetLikeRepo := infraDB.NewTweetLikeRepository(db)

	// Initialize email service
	emailHost := os.Getenv("EMAIL_HOST")
	emailPort, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	emailUsername := os.Getenv("EMAIL_USERNAME")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailService := email.NewEmailService(emailHost, emailPort, emailUsername, emailPassword)

	// Initialize webpush service
	webPushService := webpush.NewWebPushService()

	// Initialize R2 storage service
	r2Service := storage.NewR2Service()

	// メール認証機能の有効/無効を環境変数から取得（デフォルトは無効）
	enableEmailVerification := os.Getenv("ENABLE_EMAIL_VERIFICATION") == "true"

	// JWT秘密鍵を環境変数から取得（デフォルト値を設定）
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		jwtSecretKey = "meguru-secret-key-default-2024" // デフォルト値（本番では必ず環境変数を設定）
	}

	// Initialize JWT service
	jwtService := usecase.NewJWTService(jwtSecretKey)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo, passwordResetTokenRepo, emailService, jwtService)
	healthUsecase := usecase.NewHealthUsecase()
	storeUsecase := usecase.NewStoreUsecase(storeRepo, storeTokenRepo, emailService, enableEmailVerification)
	pushSubscriptionUsecase := usecase.NewPushSubscriptionUsecase(pushSubscriptionRepo, webPushService)
	flyerUsecase := usecase.NewFlyerUsecase(flyerRepo, storeRepo, productRepo)
	productUsecase := usecase.NewProductUsecase(productRepo, storeRepo, r2Service)
	newsViewUsecase := usecase.NewNewsViewUsecase(newsViewRepo, storeRepo)
	tweetUsecase := usecase.NewTweetUsecase(tweetRepo, tweetLikeRepo, storeRepo)

	// Initialize controllers
	userController := controller.NewUserController(userUsecase)
	healthController := controller.NewHealthController(healthUsecase)
	storeController := controller.NewStoreController(storeUsecase)
	pushSubscriptionController := controller.NewPushSubscriptionController(pushSubscriptionUsecase)
	flyerController := controller.NewFlyerController(flyerUsecase)
	productController := controller.NewProductController(productUsecase)
	newsViewController := controller.NewNewsViewController(newsViewUsecase)
	tweetController := controller.NewTweetController(tweetUsecase)

	// Initialize router
	r := router.NewRouter(userController, healthController, storeController, pushSubscriptionController, flyerController, productController, newsViewController, tweetController, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}