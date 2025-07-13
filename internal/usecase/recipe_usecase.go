package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/infrastructure/r2"
)

type RecipeUsecase struct {
	recipeRepo repository.RecipeRepository
	r2Service  *r2.R2Service
	projectID  string
	location   string
}

func NewRecipeUsecase(recipeRepo repository.RecipeRepository, r2Service *r2.R2Service, projectID, location string) *RecipeUsecase {
	return &RecipeUsecase{
		recipeRepo: recipeRepo,
		r2Service:  r2Service,
		projectID:  projectID,
		location:   location,
	}
}

func (u *RecipeUsecase) SuggestRecipes(ctx context.Context, ingredients []string) ([]entity.Recipe, error) {
	recipes, err := u.recipeRepo.FindRecipesByIngredients(ctx, ingredients)
	if err != nil {
		return nil, fmt.Errorf("failed to find recipes by ingredients: %w", err)
	}

	if len(recipes) > 0 {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(recipes), func(i, j int) { recipes[i], recipes[j] = recipes[j], recipes[i] })
		if len(recipes) > 3 {
			return recipes[:3], nil
		}
		return recipes, nil
	}

	client, err := genai.NewClient(ctx, u.projectID, u.location)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash-001")

	prompt := fmt.Sprintf(`以下の材料を使ったレシピを10個提案してください。

材料: %s

各レシピは以下のJSON形式で、配列として返してください。JSONオブジェクトのみを生成し、マークダウンのバッククォート("""json ... """)は含めないでください。

[
  {
    "title": "レシピ名",
    "description": "レシピの説明",
    "cooking_time": 調理時間(分),
    "servings": 何人前か,
	"cost": 予想金額(円),
	"total_calories": 総カロリー(kcal),
	"total_score": 総合評価(10点満点。調理時間、費用、カロリーを総合的に評価),
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

	cleanedJsonString := strings.TrimSpace(jsonString)
	cleanedJsonString = strings.TrimPrefix(cleanedJsonString, "```json")
	cleanedJsonString = strings.TrimSuffix(cleanedJsonString, "```")
	cleanedJsonString = strings.TrimSpace(cleanedJsonString)

	var suggestedRecipes []struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		CookingTime   int    `json:"cooking_time"`
		Servings      int    `json:"servings"`
		Cost          int    `json:"cost"`
		TotalCalories int    `json:"total_calories"`
		TotalScore    int    `json:"total_score"`
		Ingredients   []struct {
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
		// --- Image Generation ---
		imagePrompt := fmt.Sprintf("A realistic, high-quality photo of %s", sr.Title)
		imgModel := client.GenerativeModel("imagegeneration@006")
		imageResp, err := imgModel.GenerateContent(ctx, genai.Text(imagePrompt))

		if err != nil {
			log.Printf("failed to generate image for %s: %v", sr.Title, err)
		}

		var imageURL string
		if imageResp != nil && len(imageResp.Candidates) > 0 && imageResp.Candidates[0].Content != nil {
			for _, part := range imageResp.Candidates[0].Content.Parts {
				if blob, ok := part.(genai.Blob); ok {
					fileKey := fmt.Sprintf("recipe-images/%s-%d.png", strings.ReplaceAll(sr.Title, " ", "-"), time.Now().UnixNano())
					imageURL, err = u.r2Service.UploadImage(ctx, fileKey, blob.Data)
					if err != nil {
						log.Printf("failed to upload image for %s: %v", sr.Title, err)
					}
					break
				}
			}
		}

		// --- Prepare data for saving ---
		recipe := entity.Recipe{
			Title:         sr.Title,
			Description:   sr.Description,
			ImageURL:      imageURL,
			CookingTime:   sr.CookingTime,
			Servings:      sr.Servings,
			Cost:          sr.Cost,
			TotalCalories: sr.TotalCalories,
			TotalScore:    sr.TotalScore,
		}

		// --- Combine original and suggested ingredients ---
		ingredientSet := make(map[string]struct{})
		var finalIngredients []entity.Ingredient

		// Add original ingredients from the request first
		for _, name := range ingredients {
			if _, exists := ingredientSet[name]; !exists {
				ingredientSet[name] = struct{}{}
				finalIngredients = append(finalIngredients, entity.Ingredient{Name: name})
			}
		}
		// Add ingredients from the Gemini response
		for _, i := range sr.Ingredients {
			if _, exists := ingredientSet[i.Name]; !exists {
				ingredientSet[i.Name] = struct{}{}
				finalIngredients = append(finalIngredients, entity.Ingredient{Name: i.Name})
			}
		}

		var recipeSteps []entity.Step
		for _, s := range sr.Steps {
			recipeSteps = append(recipeSteps, entity.Step{StepNumber: s.StepNumber, Instruction: s.Instruction})
		}

		savedRecipe, err := u.recipeRepo.SaveRecipe(ctx, recipe, finalIngredients, recipeSteps)
		if err != nil {
			return nil, fmt.Errorf("failed to save recipe: %w", err)
		}
		savedRecipes = append(savedRecipes, savedRecipe)
	}

	if len(savedRecipes) > 3 {
		return savedRecipes[:3], nil
	}

	return savedRecipes, nil
}
