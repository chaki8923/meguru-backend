package router

import (
	"meguru-backend/internal/interface/controller"
	"meguru-backend/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(userController *controller.UserController, healthController *controller.HealthController, storeController *controller.StoreController, pushSubscriptionController *controller.PushSubscriptionController, flyerController *controller.FlyerController, flyerViewController *controller.FlyerViewController, productController *controller.ProductController, newsViewController *controller.NewsViewController, newsConsultationController *controller.NewsConsultationController, tweetController *controller.TweetController, recipeController *controller.RecipeController, jwtService *usecase.JWTService) *gin.Engine {
	r := gin.Default()

	// CORS設定
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000", 
		"http://localhost:3001",
		"http://localhost:3002",
		"http://localhost:3003",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
		"http://127.0.0.1:3002",
		"http://127.0.0.1:3003",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept", "Accept-Encoding", "X-CSRF-Token"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	// ヘルスチェックエンドポイント
	r.GET("/health", healthController.GetHealth)

	// 店舗登録エンドポイント（/store/shopRegister）
	store := r.Group("/store")
	{
		store.POST("/shopRegister", storeController.RegisterShop)
		store.GET("/verify-email", storeController.VerifyEmail) // メール認証エンドポイント
	}

	// API routes
	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/register", userController.CreateUser)
			users.POST("/login", userController.LoginUser)
			users.POST("/forgot-password", userController.ForgotPassword)
			users.POST("/reset-password", userController.ResetPassword)

			// 認証が必要なエンドポイント
			protected := users.Group("")
			protected.Use(UserAuthMiddleware(jwtService))
			{
				protected.GET("/profile", userController.GetProfile)
				protected.PUT("/profile", userController.UpdateProfile)
			}
		}
		stores := api.Group("/stores")
		{
			stores.POST("", storeController.CreateStore)
			stores.GET("", storeController.GetAllStores)
			stores.PUT("/:id", storeController.UpdateStore)
			stores.POST("/signin", storeController.SignIn)

			// 認証が必要なエンドポイント
			stores.GET("/profile", storeController.GetProfile)
			stores.PUT("/profile", storeController.UpdateProfile)
		}
		notifications := api.Group("/notifications")
		{
			notifications.POST("/subscribe", pushSubscriptionController.Subscribe)
			notifications.POST("/send", pushSubscriptionController.SendNotification)
		}
		flyer := api.Group("/flyer")
		{
			flyer.POST("/upload", flyerController.UploadFlyer)
			flyer.GET("/nearby", flyerController.GetNearbyFlyers) // 近隣店舗チラシ取得
			flyer.GET("/:store_id", flyerController.GetFlyerByStoreID)
			flyer.GET("/all/:store_id", flyerController.GetAllFlyersByStoreID)
			
			// フライヤービュー関連のエンドポイント
			if flyerViewController != nil {
				flyerViews := flyer.Group("/views")
				{
					flyerViews.POST("", flyerViewController.RecordFlyerView)                          // ビュー記録（認証不要 - ユーザー側）
					flyerViews.GET("/count/:flyer_id", flyerViewController.GetFlyerViewCount)         // ビュー数取得（店舗側）
					flyerViews.GET("/list/:flyer_id", flyerViewController.GetFlyerViewList)           // ビューリスト取得（店舗側）
				}
			}
		}

		// 商品関連のエンドポイント（認証必須）
		products := api.Group("/products")
		products.Use(AuthMiddleware()) // 店舗認証が必要
		{
			products.GET("", productController.ListStoreProducts)
			products.POST("", productController.CreateStoreProduct)
			products.GET("/:id", productController.GetStoreProduct)
			products.PUT("/:id", productController.UpdateStoreProduct)
			products.DELETE("/:id", productController.DeleteStoreProduct)
		}

		// ニュース閲覧関連のエンドポイント
		news := api.Group("/news")
		{
			news.POST("/view", newsViewController.RecordNewsView)                 // ニュース閲覧記録
			news.GET("/view-count/:news_id", newsViewController.GetNewsViewCount) // 単一ニュース閲覧数
			news.POST("/view-counts", newsViewController.GetNewsViewCounts)       // 複数ニュース閲覧数
			news.POST("/consult", newsConsultationController.ConsultNews)         // ニュース相談機能
			news.POST("/send-email", newsConsultationController.SendNewsAnalysisEmail) // ニュース分析結果メール送信
		}

		// ツイート関連のエンドポイント
		api.GET("/tweets", tweetController.ListTweets)
		tweets := api.Group("/tweets")
		tweets.Use(AuthMiddleware())
		{
			tweets.POST("", tweetController.CreateTweet)
			tweets.DELETE("/:id", tweetController.DeleteTweet)
			tweets.POST("/:id/like", tweetController.LikeTweet)
			tweets.DELETE("/:id/like", tweetController.UnlikeTweet)
		}

		// レシピ関連のエンドポイント（認証付き）
		recipes := api.Group("/recipes")
		recipes.Use(UserAuthMiddleware(jwtService))
		{
			recipes.GET("/:recipe_id", recipeController.GetRecipeDetail)
		}

		// 画像からレシピ取得エンドポイント（認証付き）
		recipesFromImage := api.Group("/recipes-from-image")
		recipesFromImage.Use(UserAuthMiddleware(jwtService))
		{
			recipesFromImage.POST("", recipeController.GetRecipesByImage)
		}

		// 保存レシピ関連のエンドポイント（認証必須）
		savedRecipes := api.Group("/saved-recipes")
		savedRecipes.Use(UserAuthMiddleware(jwtService))
		{
			savedRecipes.POST("", recipeController.SaveRecipe)
			savedRecipes.GET("", recipeController.GetSavedRecipes)
			savedRecipes.DELETE("/:recipe_id", recipeController.DeleteSavedRecipe)
		}
	}

	return r
}
