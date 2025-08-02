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
)

// Response structure to be sent to the controller
type FlyerResponse struct {
	ID         string         `json:"id"`
	StoreID    string         `json:"store_id"`
	ImageData  string         `json:"image_data"` // base64 encoded image
	FlyerData  *dto.FlyerData `json:"flyer_data"`
	CreatedAt  time.Time      `json:"created_at"`
}

type FlyerUsecase struct {
	flyerRepository repository.FlyerRepository
}

func NewFlyerUsecase(flyerRepository repository.FlyerRepository) *FlyerUsecase {
	return &FlyerUsecase{flyerRepository: flyerRepository}
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
	log.Println("プロンプト生成処理を開始します")
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
		ID:         savedFlyer.ID.String(),
		StoreID:    storeID.String(),
		ImageData:  base64.StdEncoding.EncodeToString(savedFlyer.ImageData),
		FlyerData:  &flyerData,
		CreatedAt:  savedFlyer.CreatedAt,
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
		ID:        flyer.ID.String(),
		ImageData: base64.StdEncoding.EncodeToString(flyer.ImageData),
		FlyerData: flyerData,
		CreatedAt: flyer.CreatedAt,
	}

	return response, nil
}

