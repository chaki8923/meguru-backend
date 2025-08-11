package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/infrastructure/storage"
)

// Request/Response DTOs
type CreateStoreProductRequest struct {
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	ImageURL    string `json:"image_url"`
	Status      string `json:"status"`
}

type UpdateStoreProductRequest struct {
	ProductName string `json:"product_name"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	ImageURL    string `json:"image_url"`
	Status      string `json:"status"`
}

type StoreProductResponse struct {
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Category    string    `json:"category"`
	Price       int       `json:"price"`
	Quantity    int       `json:"quantity"`
	ImageURL    string    `json:"image_url"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductUsecase struct {
	productRepository repository.ProductRepository
	storeRepository   repository.StoreRepository
	r2Service         *storage.R2Service
}

func NewProductUsecase(productRepository repository.ProductRepository, storeRepository repository.StoreRepository, r2Service *storage.R2Service) *ProductUsecase {
	return &ProductUsecase{
		productRepository: productRepository,
		storeRepository:   storeRepository,
		r2Service:         r2Service,
	}
}

// トークンから店舗IDを取得するヘルパー関数
func (u *ProductUsecase) getStoreIDFromToken(token string) (uuid.UUID, error) {
	var uuidStr string
	if strings.HasPrefix(token, "auth_token_") {
		uuidStr = strings.TrimPrefix(token, "auth_token_")
	} else if strings.HasPrefix(token, "temp_token_") {
		uuidStr = strings.TrimPrefix(token, "temp_token_")
	} else {
		return uuid.Nil, fmt.Errorf("invalid token format")
	}

	storeID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	return storeID, nil
}

// 店舗の商品を作成
func (u *ProductUsecase) CreateStoreProduct(ctx context.Context, token string, req *CreateStoreProductRequest) (*StoreProductResponse, error) {
	// トークンから店舗IDを取得
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// 商品名から商品を検索または作成
	product, err := u.productRepository.GetProductByName(ctx, req.ProductName)
	if err != nil {
		return nil, fmt.Errorf("failed to search product: %w", err)
	}

	// 商品が存在しない場合は作成
	if product == nil {
		product = &entity.Product{
			Name:     req.ProductName,
			Category: req.Category,
		}
		product, err = u.productRepository.CreateProduct(ctx, product)
		if err != nil {
			return nil, fmt.Errorf("failed to create product: %w", err)
		}
	}

	// 同じ店舗で同じ商品が既に登録されていないかチェック
	existingStoreProduct, err := u.productRepository.GetStoreProductByStoreAndProduct(ctx, storeID, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing store product: %w", err)
	}
	if existingStoreProduct != nil {
		return nil, fmt.Errorf("この商品は既に登録されています")
	}

	// 店舗商品を作成
	storeProduct := &entity.StoreProduct{
		StoreID:   storeID,
		ProductID: product.ID,
		Price:     req.Price,
		Quantity:  req.Quantity,
		ImageURL:  req.ImageURL,
		Status:    req.Status,
	}

	storeProduct, err = u.productRepository.CreateStoreProduct(ctx, storeProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to create store product: %w", err)
	}

	return &StoreProductResponse{
		ID:          storeProduct.ID.String(),
		ProductID:   product.ID.String(),
		ProductName: product.Name,
		Category:    product.Category,
		Price:       storeProduct.Price,
		Quantity:    storeProduct.Quantity,
		ImageURL:    storeProduct.ImageURL,
		Status:      storeProduct.Status,
		CreatedAt:   storeProduct.CreatedAt,
		UpdatedAt:   storeProduct.UpdatedAt,
	}, nil
}

// 画像ファイル付きで店舗の商品を作成
func (u *ProductUsecase) CreateStoreProductWithImage(ctx context.Context, token string, req *CreateStoreProductRequest, imageFile *multipart.FileHeader) (*StoreProductResponse, error) {
	// トークンから店舗IDを取得
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// 商品名から商品を検索または作成
	product, err := u.productRepository.GetProductByName(ctx, req.ProductName)
	if err != nil {
		return nil, fmt.Errorf("failed to search product: %w", err)
	}

	// 商品が存在しない場合は作成
	if product == nil {
		product = &entity.Product{
			Name:     req.ProductName,
			Category: req.Category,
		}
		product, err = u.productRepository.CreateProduct(ctx, product)
		if err != nil {
			return nil, fmt.Errorf("failed to create product: %w", err)
		}
	}

	// 同じ店舗で同じ商品が既に登録されていないかチェック
	existingStoreProduct, err := u.productRepository.GetStoreProductByStoreAndProduct(ctx, storeID, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing store product: %w", err)
	}
	if existingStoreProduct != nil {
		return nil, fmt.Errorf("この商品は既に登録されています")
	}

	// 画像をR2にアップロード
	var imageURL string
	if imageFile != nil {
		imageURL, err = u.r2Service.UploadProductImage(imageFile, storeID.String(), product.ID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to upload image: %w", err)
		}
	}

	// 店舗商品を作成
	storeProduct := &entity.StoreProduct{
		StoreID:   storeID,
		ProductID: product.ID,
		Price:     req.Price,
		Quantity:  req.Quantity,
		ImageURL:  imageURL,
		Status:    req.Status,
	}

	storeProduct, err = u.productRepository.CreateStoreProduct(ctx, storeProduct)
	if err != nil {
		// 画像のアップロードに成功していた場合は削除
		if imageURL != "" {
			_ = u.r2Service.DeleteProductImage(imageURL)
		}
		return nil, fmt.Errorf("failed to create store product: %w", err)
	}

	return &StoreProductResponse{
		ID:          storeProduct.ID.String(),
		ProductID:   product.ID.String(),
		ProductName: product.Name,
		Category:    product.Category,
		Price:       storeProduct.Price,
		Quantity:    storeProduct.Quantity,
		ImageURL:    storeProduct.ImageURL,
		Status:      storeProduct.Status,
		CreatedAt:   storeProduct.CreatedAt,
		UpdatedAt:   storeProduct.UpdatedAt,
	}, nil
}

// 店舗の商品を更新
func (u *ProductUsecase) UpdateStoreProduct(ctx context.Context, token string, storeProductID uuid.UUID, req *UpdateStoreProductRequest) (*StoreProductResponse, error) {
	// トークンから店舗IDを取得
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// 既存の店舗商品を取得
	storeProduct, err := u.productRepository.GetStoreProductByID(ctx, storeProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get store product: %w", err)
	}
	if storeProduct == nil {
		return nil, fmt.Errorf("store product not found")
	}

	// 店舗の所有権チェック
	if storeProduct.StoreID != storeID {
		return nil, fmt.Errorf("unauthorized to update this product")
	}

	// 商品名が変更された場合は商品マスターを更新
	if storeProduct.Product.Name != req.ProductName {
		product, err := u.productRepository.GetProductByName(ctx, req.ProductName)
		if err != nil {
			return nil, fmt.Errorf("failed to search product: %w", err)
		}

		// 新しい商品名の商品が存在しない場合は作成
		if product == nil {
			product = &entity.Product{
				Name:     req.ProductName,
				Category: req.Category,
			}
			product, err = u.productRepository.CreateProduct(ctx, product)
			if err != nil {
				return nil, fmt.Errorf("failed to create product: %w", err)
			}
		}

		storeProduct.ProductID = product.ID
		storeProduct.Product = product
	}

	// 店舗商品情報を更新
	storeProduct.Price = req.Price
	storeProduct.Quantity = req.Quantity
	// 画像URLは新しいリクエストに含まれている場合のみ更新
	if req.ImageURL != "" {
		storeProduct.ImageURL = req.ImageURL
	}
	storeProduct.Status = req.Status

	storeProduct, err = u.productRepository.UpdateStoreProduct(ctx, storeProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to update store product: %w", err)
	}

	return &StoreProductResponse{
		ID:          storeProduct.ID.String(),
		ProductID:   storeProduct.Product.ID.String(),
		ProductName: storeProduct.Product.Name,
		Category:    storeProduct.Product.Category,
		Price:       storeProduct.Price,
		Quantity:    storeProduct.Quantity,
		ImageURL:    storeProduct.ImageURL,
		Status:      storeProduct.Status,
		CreatedAt:   storeProduct.CreatedAt,
		UpdatedAt:   storeProduct.UpdatedAt,
	}, nil
}

// 画像ファイル付きで店舗の商品を更新
func (u *ProductUsecase) UpdateStoreProductWithImage(ctx context.Context, token string, storeProductID uuid.UUID, req *UpdateStoreProductRequest, imageFile *multipart.FileHeader) (*StoreProductResponse, error) {
	// トークンから店舗IDを取得
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// 既存の店舗商品を取得
	storeProduct, err := u.productRepository.GetStoreProductByID(ctx, storeProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get store product: %w", err)
	}
	if storeProduct == nil {
		return nil, fmt.Errorf("store product not found")
	}

	// 店舗の所有権チェック
	if storeProduct.StoreID != storeID {
		return nil, fmt.Errorf("unauthorized to update this product")
	}

	// 既存の画像URLを保存（削除用）
	oldImageURL := storeProduct.ImageURL

	// 商品名が変更された場合は商品マスターを更新
	if storeProduct.Product.Name != req.ProductName {
		product, err := u.productRepository.GetProductByName(ctx, req.ProductName)
		if err != nil {
			return nil, fmt.Errorf("failed to search product: %w", err)
		}

		// 新しい商品名の商品が存在しない場合は作成
		if product == nil {
			product = &entity.Product{
				Name:     req.ProductName,
				Category: req.Category,
			}
			product, err = u.productRepository.CreateProduct(ctx, product)
			if err != nil {
				return nil, fmt.Errorf("failed to create product: %w", err)
			}
		}

		storeProduct.ProductID = product.ID
		storeProduct.Product = product
	}

	// 新しい画像がアップロードされた場合
	if imageFile != nil {
		imageURL, err := u.r2Service.UploadProductImage(imageFile, storeID.String(), storeProduct.Product.ID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to upload image: %w", err)
		}
		storeProduct.ImageURL = imageURL
	}

	// 店舗商品情報を更新
	storeProduct.Price = req.Price
	storeProduct.Quantity = req.Quantity
	storeProduct.Status = req.Status

	storeProduct, err = u.productRepository.UpdateStoreProduct(ctx, storeProduct)
	if err != nil {
		// 新しい画像のアップロードに成功していた場合は削除
		if imageFile != nil && storeProduct.ImageURL != oldImageURL {
			_ = u.r2Service.DeleteProductImage(storeProduct.ImageURL)
		}
		return nil, fmt.Errorf("failed to update store product: %w", err)
	}

	// 古い画像を削除（新しい画像がアップロードされた場合のみ）
	if imageFile != nil && oldImageURL != "" && oldImageURL != storeProduct.ImageURL {
		_ = u.r2Service.DeleteProductImage(oldImageURL)
	}

	return &StoreProductResponse{
		ID:          storeProduct.ID.String(),
		ProductID:   storeProduct.Product.ID.String(),
		ProductName: storeProduct.Product.Name,
		Category:    storeProduct.Product.Category,
		Price:       storeProduct.Price,
		Quantity:    storeProduct.Quantity,
		ImageURL:    storeProduct.ImageURL,
		Status:      storeProduct.Status,
		CreatedAt:   storeProduct.CreatedAt,
		UpdatedAt:   storeProduct.UpdatedAt,
	}, nil
}

// 店舗の商品を削除
func (u *ProductUsecase) DeleteStoreProduct(ctx context.Context, token string, storeProductID uuid.UUID) error {
	// トークンから店舗IDを取得
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// 既存の店舗商品を取得
	storeProduct, err := u.productRepository.GetStoreProductByID(ctx, storeProductID)
	if err != nil {
		return fmt.Errorf("failed to get store product: %w", err)
	}
	if storeProduct == nil {
		return fmt.Errorf("store product not found")
	}

	// 店舗の所有権チェック
	if storeProduct.StoreID != storeID {
		return fmt.Errorf("unauthorized to delete this product")
	}

	// 商品に関連付けられた画像を保存
	imageURL := storeProduct.ImageURL

	err = u.productRepository.DeleteStoreProduct(ctx, storeProductID)
	if err != nil {
		return fmt.Errorf("failed to delete store product: %w", err)
	}

	// 画像も削除
	if imageURL != "" {
		_ = u.r2Service.DeleteProductImage(imageURL)
	}

	return nil
}

// 店舗の商品一覧を取得
func (u *ProductUsecase) ListStoreProducts(ctx context.Context, token string, includeExpired bool) ([]*StoreProductResponse, error) {
	// トークンから店舗IDを取得
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	storeProducts, err := u.productRepository.ListStoreProductsByStoreID(ctx, storeID, includeExpired)
	if err != nil {
		return nil, fmt.Errorf("failed to list store products: %w", err)
	}

	var responses []*StoreProductResponse
	for _, sp := range storeProducts {
		responses = append(responses, &StoreProductResponse{
			ID:          sp.ID.String(),
			ProductID:   sp.Product.ID.String(),
			ProductName: sp.Product.Name,
			Category:    sp.Product.Category,
			Price:       sp.Price,
			Quantity:    sp.Quantity,
			ImageURL:    sp.ImageURL,
			Status:      sp.Status,
			CreatedAt:   sp.CreatedAt,
			UpdatedAt:   sp.UpdatedAt,
		})
	}

	return responses, nil
}

// 店舗の商品を1件取得
func (u *ProductUsecase) GetStoreProduct(ctx context.Context, token string, storeProductID uuid.UUID) (*StoreProductResponse, error) {
	// トークンから店舗IDを取得
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	storeProduct, err := u.productRepository.GetStoreProductByID(ctx, storeProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get store product: %w", err)
	}
	if storeProduct == nil {
		return nil, fmt.Errorf("store product not found")
	}

	// 店舗の所有権チェック
	if storeProduct.StoreID != storeID {
		return nil, fmt.Errorf("unauthorized to access this product")
	}

	return &StoreProductResponse{
		ID:          storeProduct.ID.String(),
		ProductID:   storeProduct.Product.ID.String(),
		ProductName: storeProduct.Product.Name,
		Category:    storeProduct.Product.Category,
		Price:       storeProduct.Price,
		Quantity:    storeProduct.Quantity,
		ImageURL:    storeProduct.ImageURL,
		Status:      storeProduct.Status,
		CreatedAt:   storeProduct.CreatedAt,
		UpdatedAt:   storeProduct.UpdatedAt,
	}, nil
} 