package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"regexp"
	"strings"
	"time"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/dto"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"github.com/google/uuid"
)

// Response structure to be sent to the controller
type FlyerResponse struct {
	ID                string         `json:"id"`
	StoreID           string         `json:"store_id"`
	ImageData         string         `json:"image_data"` // base64 encoded image
	FlyerData         *dto.FlyerData `json:"flyer_data"`
	DisplayExpiryDate *time.Time     `json:"display_expiry_date,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
}

type FlyerUsecase struct {
	flyerRepository  repository.FlyerRepository
	storeRepository  repository.StoreRepository
	productRepository repository.ProductRepository
}

func NewFlyerUsecase(flyerRepository repository.FlyerRepository, storeRepository repository.StoreRepository, productRepository repository.ProductRepository) *FlyerUsecase {
	return &FlyerUsecase{
		flyerRepository:  flyerRepository,
		storeRepository:  storeRepository,
		productRepository: productRepository,
	}
}

func (u *FlyerUsecase) AnalyzeAndSaveFlyer(ctx context.Context, fileHeader *multipart.FileHeader) (*FlyerResponse, error) {
	// 1. Read the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// 2. Call Gemini API
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	log.Println("プロンプト生成処理を開始します!")
	prompt := `添付されたスーパーのチラシ画像を分析し、以下のJSON形式で情報を抽出してください。

出力形式のルール:
- JSONオブジェクトのみを生成してください。マークダウンのバッククォート("""json ... """)は含めないでください。
- すべての情報は指定されたJSON構造に従う必要があります。
- 日付は "YYYY-MM-DD" 形式で記述してください。
- 価格は数値型(integer)で設定してください。

JSON構造:
{
  "store": {
    "name": "店舗名",
    "prefecture": "都道府県",
    "city": "市区町村",
    "street": "番地"
  },
  "campaign": {
    "name": "キャンペーン名 (例: スーパー火曜祭)",
    "start_date": "開始日",
    "end_date": "終了日"
  },
  "flyer_items": [
    {
      "product": {
        "name": "商品名",
        "category": "カテゴリ"
      },
      "price_excluding_tax": 0,
      "price_including_tax": 0,
      "unit": "単位 (例: 各, 1個)",
      "restriction_note": "購入制限 (例: お一人様2点限り)"
    }
  ]
}

上記の指示に従って、JSONオブジェクトのみを出力してください。`

	resp, err := model.GenerateContent(ctx, genai.Text(prompt), genai.ImageData("png", imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// 3. Extract and parse JSON response
	var jsonString string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					jsonString += string(txt)
				}
			}
		}
	}

	// Clean up the JSON string
	re := regexp.MustCompile("(?s)```json(.*)```")
	matches := re.FindStringSubmatch(jsonString)
	if len(matches) > 1 {
		jsonString = strings.TrimSpace(matches[1])
	} else {
		jsonString = strings.TrimSpace(jsonString)
	}

	if jsonString == "" {
		return nil, fmt.Errorf("no JSON generated from Gemini")
	}

	log.Printf("Generated JSON: %s", jsonString)

	var flyerData dto.FlyerData
	if err := json.Unmarshal([]byte(jsonString), &flyerData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// 4. Save data to the database
	flyerToSave := &entity.Flyer{
		ImageData: imageData,
	}

	savedFlyer, storeID, err := u.flyerRepository.SaveFlyer(ctx, flyerToSave, &flyerData)
	if err != nil {
		return nil, fmt.Errorf("failed to save flyer data: %w", err)
	}

	// 5. Construct and return the response
	response := &FlyerResponse{
		ID:                savedFlyer.ID.String(),
		StoreID:           storeID.String(),
		ImageData:         base64.StdEncoding.EncodeToString(savedFlyer.ImageData),
		FlyerData:         &flyerData,
		DisplayExpiryDate: savedFlyer.DisplayExpiryDate,
		CreatedAt:         savedFlyer.CreatedAt,
	}

	return response, nil
}

func (u *FlyerUsecase) GetFlyerByStoreID(ctx context.Context, storeID string) (*FlyerResponse, error) {
	flyer, flyerData, err := u.flyerRepository.GetFlyerByStoreID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get flyer from repository: %w", err)
	}
	if flyer == nil {
		return nil, nil // No flyer found
	}

	response := &FlyerResponse{
		ID:                flyer.ID.String(),
		ImageData:         base64.StdEncoding.EncodeToString(flyer.ImageData),
		FlyerData:         flyerData,
		DisplayExpiryDate: flyer.DisplayExpiryDate,
		CreatedAt:         flyer.CreatedAt,
	}

	return response, nil
}

// ログイン中の店舗情報をチラシ分析結果で更新する新しいメソッド
func (u *FlyerUsecase) AnalyzeAndUpdateStoreFromFlyer(ctx context.Context, fileHeader *multipart.FileHeader, token string) (*FlyerResponse, error) {
	// 1. トークンから店舗IDを取得
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// 2. 現在の店舗情報を取得
	currentStore, err := u.storeRepository.FindByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current store: %w", err)
	}
	if currentStore == nil {
		return nil, fmt.Errorf("store not found")
	}

	// 3. チラシ画像を読み込み
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// 4. AI分析でチラシ情報を抽出
	flyerData, err := u.analyzeFlyer(ctx, imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze flyer: %w", err)
	}

	// 5. 店舗情報を更新（空でない値のみ更新）
	if flyerData.StoreInfo.Name != "" {
		currentStore.Name = flyerData.StoreInfo.Name
	}
	if flyerData.StoreInfo.Prefecture != "" {
		currentStore.Prefecture = flyerData.StoreInfo.Prefecture
	}
	if flyerData.StoreInfo.City != "" {
		currentStore.City = flyerData.StoreInfo.City
	}
	if flyerData.StoreInfo.Street != "" {
		currentStore.Street = flyerData.StoreInfo.Street
	}
	currentStore.UpdatedAt = time.Now()

	// 6. 店舗情報を更新
	if err := u.storeRepository.Update(ctx, currentStore); err != nil {
		return nil, fmt.Errorf("failed to update store: %w", err)
	}

	// 7. チラシから抽出された商品を店舗商品として自動登録
	if flyerData.FlyerItemsInfo != nil && len(flyerData.FlyerItemsInfo) > 0 {
		for _, item := range flyerData.FlyerItemsInfo {
			if item.Product.Name == "" {
				continue // 商品名が空の場合はスキップ
			}

			// 商品マスターを検索または作成
			product, err := u.productRepository.GetProductByName(ctx, item.Product.Name)
			if err != nil {
				log.Printf("Error searching product %s: %v", item.Product.Name, err)
				continue
			}

			if product == nil {
				// 商品が存在しない場合は作成
				product = &entity.Product{
					Name:     item.Product.Name,
					Category: item.Product.Category,
				}
				product, err = u.productRepository.CreateProduct(ctx, product)
				if err != nil {
					log.Printf("Error creating product %s: %v", item.Product.Name, err)
					continue
				}
			}

			// 既に店舗商品として登録されているかチェック
			existingStoreProduct, err := u.productRepository.GetStoreProductByStoreAndProduct(ctx, storeID, product.ID)
			if err != nil {
				log.Printf("Error checking existing store product %s: %v", item.Product.Name, err)
				continue
			}

			if existingStoreProduct != nil {
				// 既存の店舗商品の場合は価格と在庫を更新
				existingStoreProduct.Price = item.PriceIncludingTax
				existingStoreProduct.Status = "在庫あり" // チラシに載っているので在庫ありとする

				_, err = u.productRepository.UpdateStoreProduct(ctx, existingStoreProduct)
				if err != nil {
					log.Printf("Error updating store product %s: %v", item.Product.Name, err)
				} else {
					log.Printf("Updated store product: %s (price: %d)", item.Product.Name, item.PriceIncludingTax)
				}
			} else {
				// 新規店舗商品として登録
				storeProduct := &entity.StoreProduct{
					StoreID:   storeID,
					ProductID: product.ID,
					Price:     item.PriceIncludingTax,
					Quantity:  1, // デフォルト在庫数
					ImageURL:  "", // チラシ画像からは個別商品画像は取得できない
					Status:    "在庫あり",
				}

				_, err = u.productRepository.CreateStoreProduct(ctx, storeProduct)
				if err != nil {
					log.Printf("Error creating store product %s: %v", item.Product.Name, err)
				} else {
					log.Printf("Created new store product: %s (price: %d)", item.Product.Name, item.PriceIncludingTax)
				}
			}
		}
	}

	// 8. チラシ情報を保存
	flyerToSave := &entity.Flyer{
		ImageData: imageData,
	}

	savedFlyer, _, err := u.flyerRepository.SaveFlyerForStore(ctx, flyerToSave, flyerData, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to save flyer data: %w", err)
	}

	// 9. レスポンス作成
	response := &FlyerResponse{
		ID:                savedFlyer.ID.String(),
		StoreID:           storeID.String(),
		ImageData:         base64.StdEncoding.EncodeToString(savedFlyer.ImageData),
		FlyerData:         flyerData,
		DisplayExpiryDate: savedFlyer.DisplayExpiryDate,
		CreatedAt:         savedFlyer.CreatedAt,
	}

	return response, nil
}

// トークンからストアIDを取得するヘルパー関数（storeユースケースと同じロジック）
func (u *FlyerUsecase) getStoreIDFromToken(token string) (uuid.UUID, error) {
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

// AI分析部分を抽出したヘルパー関数
func (u *FlyerUsecase) analyzeFlyer(ctx context.Context, imageData []byte) (*dto.FlyerData, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	log.Println("プロンプト生成処理を開始します!")
	prompt := `添付されたスーパーのチラシ画像を分析し、以下のJSON形式で情報を抽出してください。

出力形式のルール:
- JSONオブジェクトのみを生成してください。マークダウンのバッククォート("""json ... """)は含めないでください。
- すべての情報は指定されたJSON構造に従う必要があります。
- 日付は "YYYY-MM-DD" 形式で記述してください。
- 価格は数値型(integer)で設定してください。

JSON構造:
{
  "store": {
    "name": "店舗名",
    "prefecture": "都道府県",
    "city": "市区町村",
    "street": "番地"
  },
  "campaign": {
    "name": "キャンペーン名 (例: スーパー火曜祭)",
    "start_date": "開始日",
    "end_date": "終了日"
  },
  "flyer_items": [
    {
      "product": {
        "name": "商品名",
        "category": "カテゴリ"
      },
      "price_excluding_tax": 0,
      "price_including_tax": 0,
      "unit": "単位 (例: 各, 1個)",
      "restriction_note": "購入制限 (例: お一人様2点限り)"
    }
  ]
}

上記の指示に従って、JSONオブジェクトのみを出力してください。`

	resp, err := model.GenerateContent(ctx, genai.Text(prompt), genai.ImageData("png", imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// JSON応答を抽出・解析
	var jsonString string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					jsonString += string(txt)
				}
			}
		}
	}

	// JSON文字列をクリーンアップ
	re := regexp.MustCompile("(?s)```json(.*)```")
	matches := re.FindStringSubmatch(jsonString)
	if len(matches) > 1 {
		jsonString = strings.TrimSpace(matches[1])
	} else {
		jsonString = strings.TrimSpace(jsonString)
	}

	if jsonString == "" {
		return nil, fmt.Errorf("no JSON generated from Gemini")
	}

	log.Printf("Generated JSON: %s", jsonString)

	var flyerData dto.FlyerData
	if err := json.Unmarshal([]byte(jsonString), &flyerData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &flyerData, nil
}

// AnalyzeAndUpdateStoreFromFlyerWithData - 追加のflyerDataを受け取るバージョン
func (u *FlyerUsecase) AnalyzeAndUpdateStoreFromFlyerWithData(ctx context.Context, fileHeader *multipart.FileHeader, token string, flyerDataUpdate *dto.FlyerData) (*FlyerResponse, error) {
	// 1. トークンから店舗IDを取得
	storeID, err := u.getStoreIDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// 2. 現在の店舗情報を取得
	currentStore, err := u.storeRepository.FindByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current store: %w", err)
	}
	if currentStore == nil {
		return nil, fmt.Errorf("store not found")
	}

	// 3. チラシ画像を読み込み
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// 4. AI分析でチラシ情報を抽出
	flyerData, err := u.analyzeFlyer(ctx, imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze flyer: %w", err)
	}

	// 5. flyerDataUpdateからDisplayExpiryDateを設定
	if flyerDataUpdate != nil && flyerDataUpdate.DisplayExpiryDate != nil {
		flyerData.DisplayExpiryDate = flyerDataUpdate.DisplayExpiryDate
	}

	// 6. 店舗情報を更新（空でない値のみ更新）
	if flyerData.StoreInfo.Name != "" {
		currentStore.Name = flyerData.StoreInfo.Name
	}
	if flyerData.StoreInfo.Prefecture != "" {
		currentStore.Prefecture = flyerData.StoreInfo.Prefecture
	}
	if flyerData.StoreInfo.City != "" {
		currentStore.City = flyerData.StoreInfo.City
	}
	if flyerData.StoreInfo.Street != "" {
		currentStore.Street = flyerData.StoreInfo.Street
	}
	currentStore.UpdatedAt = time.Now()

	// 7. 店舗情報を更新
	if err := u.storeRepository.Update(ctx, currentStore); err != nil {
		return nil, fmt.Errorf("failed to update store: %w", err)
	}

	// 8. チラシから抽出された商品を店舗商品として自動登録
	if flyerData.FlyerItemsInfo != nil && len(flyerData.FlyerItemsInfo) > 0 {
		for _, item := range flyerData.FlyerItemsInfo {
			if item.Product.Name == "" {
				continue // 商品名が空の場合はスキップ
			}

			// 商品マスターを検索または作成
			product, err := u.productRepository.GetProductByName(ctx, item.Product.Name)
			if err != nil {
				log.Printf("Error searching product %s: %v", item.Product.Name, err)
				continue
			}

			if product == nil {
				// 商品が存在しない場合は作成
				product = &entity.Product{
					Name:     item.Product.Name,
					Category: item.Product.Category,
				}
				product, err = u.productRepository.CreateProduct(ctx, product)
				if err != nil {
					log.Printf("Error creating product %s: %v", item.Product.Name, err)
					continue
				}
			}

			// 既に店舗商品として登録されているかチェック
			existingStoreProduct, err := u.productRepository.GetStoreProductByStoreAndProduct(ctx, storeID, product.ID)
			if err != nil {
				log.Printf("Error checking existing store product %s: %v", item.Product.Name, err)
				continue
			}

			if existingStoreProduct != nil {
				// 既存の店舗商品の場合は価格と在庫を更新
				existingStoreProduct.Price = item.PriceIncludingTax
				existingStoreProduct.Status = "在庫あり" // チラシに載っているので在庫ありとする

				_, err = u.productRepository.UpdateStoreProduct(ctx, existingStoreProduct)
				if err != nil {
					log.Printf("Error updating store product %s: %v", item.Product.Name, err)
				} else {
					log.Printf("Updated store product: %s (price: %d)", item.Product.Name, item.PriceIncludingTax)
				}
			} else {
				// 新規店舗商品として登録
				storeProduct := &entity.StoreProduct{
					StoreID:   storeID,
					ProductID: product.ID,
					Price:     item.PriceIncludingTax,
					Quantity:  1, // デフォルト在庫数
					ImageURL:  "", // チラシ画像からは個別商品画像は取得できない
					Status:    "在庫あり",
				}

				_, err = u.productRepository.CreateStoreProduct(ctx, storeProduct)
				if err != nil {
					log.Printf("Error creating store product %s: %v", item.Product.Name, err)
				} else {
					log.Printf("Created new store product: %s (price: %d)", item.Product.Name, item.PriceIncludingTax)
				}
			}
		}
	}

	// 9. チラシ情報を保存
	flyerToSave := &entity.Flyer{
		ImageData: imageData,
	}

	savedFlyer, _, err := u.flyerRepository.SaveFlyerForStore(ctx, flyerToSave, flyerData, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to save flyer data: %w", err)
	}

	// 10. レスポンス作成
	response := &FlyerResponse{
		ID:                savedFlyer.ID.String(),
		StoreID:           storeID.String(),
		ImageData:         base64.StdEncoding.EncodeToString(savedFlyer.ImageData),
		FlyerData:         flyerData,
		DisplayExpiryDate: savedFlyer.DisplayExpiryDate,
		CreatedAt:         savedFlyer.CreatedAt,
	}

	return response, nil
}

// GetAllFlyersByStoreID retrieves all flyers for a specific store ID
func (u *FlyerUsecase) GetAllFlyersByStoreID(ctx context.Context, storeID string) ([]*FlyerResponse, error) {
	log.Printf("FlyerUsecase: Getting all flyers for storeID: %s", storeID)

	// Get all flyers from repository
	flyers, flyerDataList, err := u.flyerRepository.GetAllFlyersByStoreID(ctx, storeID)
	if err != nil {
		log.Printf("FlyerUsecase: Error getting all flyers from repository: %v", err)
		return nil, fmt.Errorf("failed to get all flyers from repository: %w", err)
	}

	if len(flyers) == 0 {
		log.Printf("FlyerUsecase: No flyers found for storeID: %s", storeID)
		return []*FlyerResponse{}, nil // Return empty slice instead of nil
	}

	// Convert to response format
	responses := make([]*FlyerResponse, len(flyers))
	for i, flyer := range flyers {
		imageDataStr := base64.StdEncoding.EncodeToString(flyer.ImageData)
		
		responses[i] = &FlyerResponse{
			ID:                flyer.ID.String(),
			StoreID:           storeID,
			ImageData:         imageDataStr,
			FlyerData:         flyerDataList[i],
			DisplayExpiryDate: flyer.DisplayExpiryDate,
			CreatedAt:         flyer.CreatedAt,
		}
	}

	log.Printf("FlyerUsecase: Successfully retrieved %d flyers for storeID: %s", len(responses), storeID)
	return responses, nil
}

// GetNearbyFlyers 近隣店舗のチラシを取得
func (u *FlyerUsecase) GetNearbyFlyers(ctx context.Context, city string, limit int) ([]*FlyerResponse, error) {
	log.Printf("FlyerUsecase: Getting nearby flyers for city: %s, limit: %d", city, limit)

	// 近隣店舗のチラシを取得
	flyers, flyerDataList, err := u.flyerRepository.GetNearbyFlyers(ctx, city, limit)
	if err != nil {
		log.Printf("Error getting nearby flyers: %v", err)
		return nil, fmt.Errorf("failed to get nearby flyers: %w", err)
	}

	if len(flyers) == 0 {
		log.Printf("No nearby flyers found for city: %s", city)
		return []*FlyerResponse{}, nil
	}

	responses := make([]*FlyerResponse, len(flyers))
	for i, flyer := range flyers {
		// 画像データをbase64エンコード
		imageDataStr := base64.StdEncoding.EncodeToString(flyer.ImageData)

		responses[i] = &FlyerResponse{
			ID:                flyer.ID.String(),
			StoreID:           flyerDataList[i].StoreInfo.ID,
			ImageData:         imageDataStr,
			FlyerData:         flyerDataList[i],
			DisplayExpiryDate: flyer.DisplayExpiryDate,
			CreatedAt:         flyer.CreatedAt,
		}
	}

	log.Printf("FlyerUsecase: Successfully retrieved %d nearby flyers for city: %s", len(responses), city)
	return responses, nil
}

