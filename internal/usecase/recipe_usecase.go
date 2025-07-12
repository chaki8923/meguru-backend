package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type RecipeUsecase struct {
	recipeRepo repository.RecipeRepository
	geminiAPIKey string
}

func NewRecipeUsecase(recipeRepo repository.RecipeRepository, geminiAPIKey string) *RecipeUsecase {
	return &RecipeUsecase{
		recipeRepo: recipeRepo,
		geminiAPIKey: geminiAPIKey,
	}
}

func (u *RecipeUsecase) SuggestRecipes(ctx context.Context, ingredients []string) ([]entity.Recipe, error) {
	recipes, err := u.recipeRepo.FindRecipesByIngredients(ctx, ingredients)
	if err != nil {
		return nil, fmt.Errorf("failed to find recipes by ingredients: %w", err)
	}

	if len(recipes) > 0 {
		return recipes, nil
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(u.geminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	prompt := fmt.Sprintf(`以下の材料を使ったレシピを3つ提案してください。

材料: %s

各レシピは以下のJSON形式で、配列として返してください。JSONオブジェクトのみを生成し、マークダウンのバッククォート("""json ... """)は含めないでください。

[
  {
    "title": "レシピ名",
    "description": "レシピの説明",
    "cooking_time": 調理時間(分),
    "servings": 何人前か,
    "ingredients": [
      {
        "name": "材料名",
        "quantity": "分量"
      }
    ],
    "steps": [
      {
        "step_number": 1,
        "instruction": "手順1"
      },
      {
        "step_number": 2,
        "instruction": "手順2"
      }
    ]
  }
]`, strings.Join(ingredients, ", "))

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

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

	log.Printf("Generated JSON: %s", jsonString)

	// Gemini APIからのレスポンスからMarkdownのバッククォートを削除
	cleanedJsonString := strings.TrimSpace(jsonString)
	cleanedJsonString = strings.TrimPrefix(cleanedJsonString, "```json")
	cleanedJsonString = strings.TrimSuffix(cleanedJsonString, "```")
	cleanedJsonString = strings.TrimSpace(cleanedJsonString)

	var suggestedRecipes []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CookingTime int    `json:"cooking_time"`
		Servings    int    `json:"servings"`
		Ingredients []struct {
			Name     string `json:"name"`
			Quantity string `json:"quantity"`
		} `json:"ingredients"`
		Steps []struct {
			StepNumber  int    `json:"step_number"`
			Instruction string `json:"instruction"`
		} `json:"steps"`
	}

	if err := json.Unmarshal([]byte(cleanedJsonString), &suggestedRecipes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	var savedRecipes []entity.Recipe
	for _, sr := range suggestedRecipes {
		recipe := entity.Recipe{
			Title:       sr.Title,
			Description: sr.Description,
			CookingTime: sr.CookingTime,
			Servings:    sr.Servings,
		}
		var recipeIngredients []entity.Ingredient
		for _, i := range sr.Ingredients {
			recipeIngredients = append(recipeIngredients, entity.Ingredient{Name: i.Name})
		}
		var recipeSteps []entity.Step
		for _, s := range sr.Steps {
			recipeSteps = append(recipeSteps, entity.Step{StepNumber: s.StepNumber, Instruction: s.Instruction})
		}

		savedRecipe, err := u.recipeRepo.SaveRecipe(ctx, recipe, recipeIngredients, recipeSteps)
		if err != nil {
			return nil, fmt.Errorf("failed to save recipe: %w", err)
		}
		savedRecipes = append(savedRecipes, savedRecipe)
	}

	return savedRecipes, nil
}
