package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/infrastructure/r2"
)

type imageGenerationJob struct {
	recipeID    uint64
	recipeTitle string
}

var (
	imageJobQueue = make(chan imageGenerationJob, 100) // Buffer for 100 jobs
	startWorkerOnce sync.Once
)

type RecipeUsecase struct {
	recipeRepo      repository.RecipeRepository
	r2Service       *r2.R2Service
	projectID       string
	location        string
	genaiClient     *genai.Client
}

func NewRecipeUsecase(recipeRepo repository.RecipeRepository, r2Service *r2.R2Service, projectID, location string) *RecipeUsecase {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		log.Fatalf("Failed to create genai client: %v", err)
	}

	u := &RecipeUsecase{
		recipeRepo:      recipeRepo,
		r2Service:       r2Service,
		projectID:       projectID,
		location:        location,
		genaiClient:     client,
	}

	// Start the worker only once
	startWorkerOnce.Do(func() {
		go u.imageGenerationWorker()
	})

	return u
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

	model := u.genaiClient.GenerativeModel("gemini-2.0-flash-001")

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
		// --- Prepare data for saving ---
		recipe := entity.Recipe{
			Title:         sr.Title,
			Description:   sr.Description,
			ImageURL:      "", // Initially empty
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

		// Add image generation job to the queue
		imageJobQueue <- imageGenerationJob{recipeID: savedRecipe.ID, recipeTitle: savedRecipe.Title}

		savedRecipes = append(savedRecipes, savedRecipe)
	}

	if len(savedRecipes) > 3 {
		return savedRecipes[:3], nil
	}

	return savedRecipes, nil
}

func (u *RecipeUsecase) imageGenerationWorker() {
	ctx := context.Background()
	imgModel := u.genaiClient.GenerativeModel("imagegeneration@006")

	for job := range imageJobQueue {
		log.Printf("[Worker] Start processing job for recipe ID %d: %s", job.recipeID, job.recipeTitle)

		var err error
		var imageURL string

		imagePrompt := fmt.Sprintf("A realistic, high-quality photo of %s", job.recipeTitle)
		imageResp, err := imgModel.GenerateContent(ctx, genai.Text(imagePrompt))
		if err != nil {
			log.Printf("[Worker] Failed to generate image for recipe ID %d: %v", job.recipeID, err)
		} else if imageResp != nil && len(imageResp.Candidates) > 0 && imageResp.Candidates[0].Content != nil {
			for _, part := range imageResp.Candidates[0].Content.Parts {
				if blob, ok := part.(genai.Blob); ok {
					fileKey := fmt.Sprintf("recipe-images/%s-%d.png", strings.ReplaceAll(job.recipeTitle, " ", "-"), time.Now().UnixNano())
					imageURL, err = u.r2Service.UploadImage(ctx, fileKey, blob.Data)
					if err != nil {
						log.Printf("[Worker] Failed to upload image for recipe ID %d: %v", job.recipeID, err)
					} else {
						log.Printf("[Worker] Image generated and uploaded for recipe ID %d. URL: %s", job.recipeID, imageURL)
						if err := u.recipeRepo.UpdateRecipeImageURL(ctx, job.recipeID, imageURL); err != nil {
							log.Printf("[Worker] Failed to update image URL for recipe ID %d: %v", job.recipeID, err)
						} else {
							log.Printf("[Worker] Successfully updated image URL for recipe ID %d", job.recipeID)
						}
					}
					break
				}
			}
		} else {
			log.Printf("[Worker] No image data found for recipe ID %d", job.recipeID)
		}

		// Wait for a moment after each job to avoid hitting the rate limit
		log.Printf("[Worker] Waiting for 5 seconds before next job...")
		time.Sleep(5 * time.Second)
	}
}
