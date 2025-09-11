package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type OpenAIService struct {
	apiKey string
	client *http.Client
}

type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type ContentItem struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
	Error   *Error   `json:"error,omitempty"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Error struct {
	Message string `json:"message"`
}

func NewOpenAIService(apiKey string) *OpenAIService {
	if apiKey == "" || apiKey == "dummy-key" {
		log.Println("Warning: OPENAI_API_KEY not set, OpenAI features will be limited")
	}

	return &OpenAIService{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (s *OpenAIService) GetIngredientsFromImage(imageBase64 string) (string, error) {
	requestBody := OpenAIRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role: "user",
				Content: []ContentItem{
					{
						Type: "text",
						Text: "画像に写っている食材と調味料の名前のみを、カンマ区切りで列挙してください。説明や文章は不要です。例：鶏肉,玉ねぎ,卵,醤油,みりん",
					},
					{
						Type: "image_url",
						ImageURL: &ImageURL{
							URL: fmt.Sprintf("data:image/jpeg;base64,%s", imageBase64),
						},
					},
				},
			},
		},
		MaxTokens: 500,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	// OpenAI APIのレスポンスは文字列として返される
	if content, ok := openAIResp.Choices[0].Message.Content.(string); ok {
		return content, nil
	}
	return "", fmt.Errorf("unexpected response format")
}

// GetRecommendedRecipes は、レシピ名のリストから AIがおすすめする10件を選択
func (s *OpenAIService) GetRecommendedRecipes(recipeNames []string) ([]string, error) {
	if s.apiKey == "" || s.apiKey == "dummy-key" {
		// APIキーが設定されていない場合は最初の10件を返す
		if len(recipeNames) > 10 {
			return recipeNames[:10], nil
		}
		return recipeNames, nil
	}

	// レシピ名を番号付きで整理
	var numberedRecipes []string
	for i, name := range recipeNames {
		numberedRecipes = append(numberedRecipes, fmt.Sprintf("%d. %s", i+1, name))
	}

	prompt := fmt.Sprintf(`以下のレシピの中から、多様性を重視して10件を選んでください。以下の観点を考慮して、バランスの良い組み合わせを選んでください：

- 料理の種類の多様性（和食、洋食、中華、エスニックなど）
- 調理方法の多様性（炒める、煮る、焼く、蒸すなど）
- 味付けの多様性（甘い、辛い、酸っぱい、しょっぱいなど）
- 食材の多様性（肉、魚、野菜、豆類など）
- 時間帯の多様性（朝食、昼食、夕食、おやつなど）

カロリー、調理時間、価格の順序は無視して、純粋に料理として魅力的で多様な組み合わせを選んでください。

レシピ一覧：
%s

回答は、選択したレシピ名のみを番号なしで、1行に1つずつ記載してください。説明や理由は不要です。最大10件まで選択してください。

例：
豚の生姜焼き
野菜炒め
味噌汁`, strings.Join(numberedRecipes, "\n"))

	requestBody := OpenAIRequest{
		Model: "gpt-4o",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 500,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	// OpenAI APIのレスポンスから推奨レシピ名を抽出
	var recommendedNames []string
	if content, ok := openAIResp.Choices[0].Message.Content.(string); ok {
		lines := strings.Split(strings.TrimSpace(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				recommendedNames = append(recommendedNames, line)
			}
		}
	}

	return recommendedNames, nil
}
