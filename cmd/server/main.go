package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	infraDB "meguru-backend/internal/infrastructure/database"
	"meguru-backend/internal/infrastructure/router"
	"meguru-backend/internal/interface/controller"
	"meguru-backend/internal/usecase"
	"meguru-backend/pkg/database"

	"github.com/joho/godotenv" // ★★★ ここのタイプミスを修正 ★★★
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

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
	userUsecase := usecase.NewUserUsecase(userRepo)
	userController := controller.NewUserController(userUsecase)
	r := router.NewRouter(userController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}