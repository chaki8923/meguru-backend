package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"meguru-backend/internal/usecase"
)

type ProductController struct {
	productUsecase *usecase.ProductUsecase
}

func NewProductController(productUsecase *usecase.ProductUsecase) *ProductController {
	return &ProductController{
		productUsecase: productUsecase,
	}
}

// Authorizationヘッダーからトークンを取得するヘルパー関数
func (c *ProductController) getTokenFromHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	log.Printf("Received Authorization header: '%s'", authHeader)
	
	if authHeader == "" {
		log.Println("Authorization header is empty")
		return "", gin.Error{Err: nil, Type: gin.ErrorTypePublic, Meta: "authorization header is required"}
	}

	// "Bearer "プレフィックスを確認して削除
	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("Authorization header does not start with 'Bearer '")
		return "", gin.Error{Err: nil, Type: gin.ErrorTypePublic, Meta: "invalid authorization format"}
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	log.Printf("Extracted token: '%s'", token)
	
	if token == "" {
		log.Println("Token is empty after removing Bearer prefix")
		return "", gin.Error{Err: nil, Type: gin.ErrorTypePublic, Meta: "token is required"}
	}

	return token, nil
}

// 商品一覧取得
func (c *ProductController) ListStoreProducts(ctx *gin.Context) {
	log.Println("ListStoreProducts controller is called")
	
	token, err := c.getTokenFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	includeExpiredStr := ctx.Query("include_expired")
	includeExpired, err := strconv.ParseBool(includeExpiredStr)
	if err != nil {
		includeExpired = false // デフォルトはfalse
	}

	products, err := c.productUsecase.ListStoreProducts(ctx.Request.Context(), token, includeExpired)
	if err != nil {
		log.Printf("Error in product usecase: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Successfully retrieved %d products", len(products))
	ctx.JSON(http.StatusOK, gin.H{"data": products})
}

// 商品作成（画像アップロード対応）
func (c *ProductController) CreateStoreProduct(ctx *gin.Context) {
	log.Println("CreateStoreProduct controller is called")
	
	token, err := c.getTokenFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// マルチパートフォームデータから商品情報を取得
	productName := ctx.PostForm("product_name")
	category := ctx.PostForm("category")
	priceStr := ctx.PostForm("price")
	quantityStr := ctx.PostForm("quantity")
	status := ctx.PostForm("status")

	if productName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "商品名が必要です"})
		return
	}

	// 価格と在庫数を数値に変換
	price := 0
	if priceStr != "" {
		if p, err := strconv.Atoi(priceStr); err == nil {
			price = p
		}
	}

	quantity := 0
	if quantityStr != "" {
		if q, err := strconv.Atoi(quantityStr); err == nil {
			quantity = q
		}
	}

	if status == "" {
		status = "在庫あり"
	}

	req := &usecase.CreateStoreProductRequest{
		ProductName: productName,
		Category:    category,
		Price:       price,
		Quantity:    quantity,
		ImageURL:    "", // R2にアップロード後に設定されるので空文字
		Status:      status,
	}

	// 画像ファイルを取得（オプション）
	imageFile, err := ctx.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("Error getting image file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルの取得に失敗しました"})
		return
	}

	log.Printf("Creating product: %+v", req)
	
	// 画像ファイルがある場合は画像付き作成、ない場合は通常作成
	var product *usecase.StoreProductResponse
	if imageFile != nil {
		product, err = c.productUsecase.CreateStoreProductWithImage(ctx.Request.Context(), token, req, imageFile)
	} else {
		product, err = c.productUsecase.CreateStoreProduct(ctx.Request.Context(), token, req)
	}

	if err != nil {
		log.Printf("Error in product usecase: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Successfully created product: %s", product.ID)
	ctx.JSON(http.StatusCreated, gin.H{"data": product})
}

// 商品詳細取得
func (c *ProductController) GetStoreProduct(ctx *gin.Context) {
	log.Println("GetStoreProduct controller is called")
	
	token, err := c.getTokenFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		log.Printf("Invalid product ID: %s", productIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効な商品IDです"})
		return
	}

	log.Printf("Getting product: %s", productID)
	
	product, err := c.productUsecase.GetStoreProduct(ctx.Request.Context(), token, productID)
	if err != nil {
		log.Printf("Error in product usecase: %v", err)
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "商品が見つかりません"})
		} else if strings.Contains(err.Error(), "unauthorized") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "この商品にアクセスする権限がありません"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("Successfully retrieved product: %s", product.ID)
	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

// 商品更新（画像アップロード対応）
func (c *ProductController) UpdateStoreProduct(ctx *gin.Context) {
	log.Println("UpdateStoreProduct controller is called")
	
	token, err := c.getTokenFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		log.Printf("Invalid product ID: %s", productIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効な商品IDです"})
		return
	}

	// マルチパートフォームデータから商品情報を取得
	productName := ctx.PostForm("product_name")
	category := ctx.PostForm("category")
	priceStr := ctx.PostForm("price")
	quantityStr := ctx.PostForm("quantity")
	status := ctx.PostForm("status")

	if productName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "商品名が必要です"})
		return
	}

	// 価格と在庫数を数値に変換
	price := 0
	if priceStr != "" {
		if p, err := strconv.Atoi(priceStr); err == nil {
			price = p
		}
	}

	quantity := 0
	if quantityStr != "" {
		if q, err := strconv.Atoi(quantityStr); err == nil {
			quantity = q
		}
	}

	if status == "" {
		status = "在庫あり"
	}

	req := &usecase.UpdateStoreProductRequest{
		ProductName: productName,
		Category:    category,
		Price:       price,
		Quantity:    quantity,
		ImageURL:    "", // R2にアップロード後に設定されるので空文字
		Status:      status,
	}

	// 画像ファイルを取得（オプション）
	imageFile, err := ctx.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("Error getting image file: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "画像ファイルの取得に失敗しました"})
		return
	}

	log.Printf("Updating product %s: %+v", productID, req)
	
	// 画像ファイルがある場合は画像付き更新、ない場合は通常更新
	var product *usecase.StoreProductResponse
	if imageFile != nil {
		product, err = c.productUsecase.UpdateStoreProductWithImage(ctx.Request.Context(), token, productID, req, imageFile)
	} else {
		product, err = c.productUsecase.UpdateStoreProduct(ctx.Request.Context(), token, productID, req)
	}

	if err != nil {
		log.Printf("Error in product usecase: %v", err)
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "商品が見つかりません"})
		} else if strings.Contains(err.Error(), "unauthorized") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "この商品を更新する権限がありません"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("Successfully updated product: %s", product.ID)
	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

// 商品削除
func (c *ProductController) DeleteStoreProduct(ctx *gin.Context) {
	log.Println("DeleteStoreProduct controller is called")
	
	token, err := c.getTokenFromHeader(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productIDStr := ctx.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		log.Printf("Invalid product ID: %s", productIDStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効な商品IDです"})
		return
	}

	log.Printf("Deleting product: %s", productID)
	
	err = c.productUsecase.DeleteStoreProduct(ctx.Request.Context(), token, productID)
	if err != nil {
		log.Printf("Error in product usecase: %v", err)
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "商品が見つかりません"})
		} else if strings.Contains(err.Error(), "unauthorized") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "この商品を削除する権限がありません"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("Successfully deleted product: %s", productID)
	ctx.JSON(http.StatusOK, gin.H{"message": "商品を削除しました"})
} 