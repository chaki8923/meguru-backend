package main

import (
    "log"
    "net/http"
    "os"
    "time"
    "github.com/gin-gonic/gin"
)

func main() {
    // 環境変数をログ出力（デバッグ用）
    log.Printf("PORT: %s", os.Getenv("PORT"))
    log.Printf("DB_HOST: %s", os.Getenv("DB_HOST"))
    log.Printf("DB_PORT: %s", os.Getenv("DB_PORT"))
    log.Printf("DB_USER: %s", os.Getenv("DB_USER"))
    log.Printf("DB_NAME: %s", os.Getenv("DB_NAME"))
    log.Printf("GIN_MODE: %s", os.Getenv("GIN_MODE"))

    // Ginルーターを初期化
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())

    // ヘルスチェックエンドポイント
    r.GET("/health", func(c *gin.Context) {
        log.Printf("Health check endpoint accessed from %s", c.ClientIP())
        
        now := time.Now()
        
        // 基本的なヘルスチェック応答
        healthResponse := gin.H{
            "status":    "okだよ",
            "service":   "meguru-backend",
            "timestamp": now.Unix(),
            "time":      now.Format(time.RFC3339),
            "uptime":    "running",
        }
        
        // 環境変数の確認（データベース設定があるかチェック）
        if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
            healthResponse["database"] = gin.H{
                "configured": true,
                "host":       dbHost,
                "status":     "connection_not_tested", // 実際のDB接続テストは後で実装
            }
        } else {
            healthResponse["database"] = gin.H{
                "configured": false,
                "status":     "not_configured",
            }
        }
        
        c.JSON(http.StatusOK, healthResponse)
    })

    // 基本的な動作確認エンドポイント
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Meguru Backend API",
            "status":  "running",
        })
    })

    // ダミーAPIエンドポイント
    api := r.Group("/api/v1")
    {
        api.GET("/users", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "message": "Users endpoint",
                "data":    []interface{}{},
            })
        })
        
        api.POST("/users", func(c *gin.Context) {
            c.JSON(http.StatusCreated, gin.H{
                "message": "User created (mock)",
            })
        })
    }

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    address := "0.0.0.0:" + port
    log.Printf("Server starting on address %s", address)
    log.Printf("Health check endpoint: http://%s/health", address)
    
    if err := r.Run(address); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}