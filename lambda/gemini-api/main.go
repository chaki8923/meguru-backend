package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiRequest struct {
	ImageData  string `json:"image_data"`  // Base64 encoded image
	Operation  string `json:"operation"`   // "validate" or "analyze"
	Prompt     string `json:"prompt"`      // Custom prompt if needed
}

type GeminiResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data"`    // JSON response from Gemini
	Error   string `json:"error"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %s", request.Body)
	
	// Parse request body
	var geminiReq GeminiRequest
	if err := json.Unmarshal([]byte(request.Body), &geminiReq); err != nil {
		log.Printf("Failed to parse request: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"success": false, "error": "Invalid request format"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Get Gemini API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("GEMINI_API_KEY not set")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"success": false, "error": "GEMINI_API_KEY not configured"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Create context with timeout
	geminiCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Call Gemini API
	result, err := callGeminiAPI(geminiCtx, apiKey, geminiReq)
	if err != nil {
		log.Printf("Gemini API call failed: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"success": false, "error": "%s"}`, err.Error()),
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
		}, nil
	}

	// Return success response
	response := GeminiResponse{
		Success: true,
		Data:    result,
	}

	responseBody, _ := json.Marshal(response)
	
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func callGeminiAPI(ctx context.Context, apiKey string, req GeminiRequest) (string, error) {
	log.Println("Creating Gemini client...")
	
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	
	// Prepare the prompt based on operation
	var prompt string
	switch req.Operation {
	case "validate":
		prompt = `添付された画像を分析して、これがスーパーマーケットのチラシ（広告）画像かどうかを判定してください。

以下のJSON形式で結果を返してください：
{
  "is_valid_flyer": true/false,
  "confidence": 0.0-1.0,
  "reason": "判定理由"
}`
	case "analyze":
		prompt = `添付されたスーパーのチラシ画像を分析し、以下のJSON形式で情報を抽出してください。

出力形式のルール:
- 必ずJSONのみを出力してください（説明文は含めない）
- 価格は数値のみ（単位記号は除く）
- 商品名は正確に抽出してください

{
  "campaign": {
    "name": "キャンペーン名",
    "start_date": "YYYY-MM-DD",
    "end_date": "YYYY-MM-DD"
  },
  "store": {
    "name": "店舗名",
    "city": "市区町村名",
    "prefecture": "都道府県名",
    "street": "番地・通り名"
  },
  "flyer_items": [
    {
      "product": {
        "name": "商品名",
        "category": "カテゴリ"
      },
      "price_excluding_tax": 価格（税抜き、数値のみ）,
      "price_including_tax": 価格（税込み、数値のみ）,
      "unit": "単位（例：1個、100g）",
      "restriction_note": "制限事項（例：お一人様3点まで）"
    }
  ]
}`
	default:
		// Use custom prompt if provided
		if req.Prompt != "" {
			prompt = req.Prompt
		} else {
			return "", fmt.Errorf("unknown operation: %s", req.Operation)
		}
	}

	// Convert base64 image data to blob
	imageBlob := genai.ImageData("image/jpeg", req.ImageData)
	
	log.Println("Sending request to Gemini API...")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt), imageBlob)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract text from response
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts in response")
	}

	responseText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	log.Printf("Gemini API response length: %d", len(responseText))
	
	return responseText, nil
}

func main() {
	lambda.Start(handleRequest)
}
