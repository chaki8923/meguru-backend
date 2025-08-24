package usecase

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"meguru-backend/internal/usecase/dto"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type NewsConsultationUsecase struct{}

func NewNewsConsultationUsecase() *NewsConsultationUsecase {
	return &NewsConsultationUsecase{}
}

// ConsultNews ニュース内容からスーパーの対応策をAIに分析させる
func (u *NewsConsultationUsecase) ConsultNews(ctx context.Context, request *dto.NewsConsultationRequest) (*dto.NewsConsultationResponse, error) {
	log.Printf("Starting news consultation for URL: %s", request.NewsURL)

	// AIに直接URLを渡して分析
	analysisResult, recommendations, err := u.analyzeNewsWithAI(ctx, request.NewsTitle, request.NewsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze news with AI: %w", err)
	}

	// レスポンスを作成
	response := &dto.NewsConsultationResponse{
		NewsID:          request.NewsID,
		NewsTitle:       request.NewsTitle,
		NewsURL:         request.NewsURL,
		AnalysisResult:  analysisResult,
		Recommendations: recommendations,
		CreatedAt:       time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	log.Printf("News consultation completed successfully for news ID: %s", request.NewsID)
	return response, nil
}

// analyzeNewsWithAI AIを使ってニュース内容を分析
func (u *NewsConsultationUsecase) analyzeNewsWithAI(ctx context.Context, title, newsURL string) (string, string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", "", fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to create genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro")

	prompt := fmt.Sprintf(`
あなたはスーパーマーケット経営のコンサルタントです。以下のニュース記事のURLにアクセスして内容を読み取り、スーパーマーケットがどのような対応を取るべきかを提案してください。

ニュースタイトル: %s
ニュースURL: %s

以下の手順で分析してください：
1. 上記URLにアクセスして記事の内容を読み取ってください
2. 記事の内容を理解し、食品小売業界への影響を分析してください
3. スーパーマーケット経営者向けの具体的な対応策を提案してください

以下の形式で回答してください：

【記事分析】
読み取った記事の内容が食品小売業界、特にスーパーマーケットに与える影響について分析してください。

【推奨対応策】
1. 短期的対応（1-3ヶ月以内）
2. 中期的対応（3-12ヶ月以内）  
3. 長期的対応（1年以上）

各対応策について、具体的で実行可能な施策を提案してください。コスト効果、実装の難易度、期待される効果も考慮に入れてください。

回答は日本語で、実際のスーパーマーケット経営者が読んで参考になるよう、具体的で実用的な内容にしてください。
`, title, newsURL)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate content: %w", err)
	}

	// レスポンスからテキストを抽出
	var fullResponse string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					fullResponse += string(txt)
				}
			}
		}
	}

	if fullResponse == "" {
		return "", "", fmt.Errorf("no response generated from AI")
	}

	// レスポンスを分析部分と推奨対応策部分に分割
	analysisResult, recommendations := u.parseAIResponse(fullResponse)

	log.Printf("AI analysis completed. Response length: %d characters", len(fullResponse))
	return analysisResult, recommendations, nil
}

// parseAIResponse AI応答を分析結果と推奨対応策に分割
func (u *NewsConsultationUsecase) parseAIResponse(response string) (string, string) {
	// 【推奨対応策】部分を境界として分割
	parts := strings.Split(response, "【推奨対応策】")
	
	analysisResult := strings.TrimSpace(parts[0])
	recommendations := ""
	
	if len(parts) > 1 {
		recommendations = "【推奨対応策】" + strings.TrimSpace(parts[1])
	}
	
	// 分析結果から【記事分析】タイトルを除去
	analysisResult = strings.ReplaceAll(analysisResult, "【記事分析】", "")
	analysisResult = strings.TrimSpace(analysisResult)
	
	return analysisResult, recommendations
}
