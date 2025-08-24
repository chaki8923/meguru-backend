package usecase

import (
	"context"
	"errors"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/infrastructure/service"
	dto "meguru-backend/internal/usecase/dto/recipes"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RecipeUsecase struct {
	recipeRepo    repository.RecipeRepository
	openAIService *service.OpenAIService
}

func NewRecipeUsecase(recipeRepo repository.RecipeRepository, openAIService *service.OpenAIService) *RecipeUsecase {
	return &RecipeUsecase{
		recipeRepo:    recipeRepo,
		openAIService: openAIService,
	}
}

func (u *RecipeUsecase) GetRecipeDetail(ctx context.Context, recipeID string) (*dto.GetRecipeDetailResponse, error) {
	// レシピ基本情報を取得
	recipe, err := u.recipeRepo.GetRecipeByID(ctx, recipeID)
	if err != nil {
		return nil, err
	}
	if recipe == nil {
		return nil, nil
	}

	// 材料情報を取得
	ingredients, err := u.recipeRepo.GetRecipeIngredients(ctx, recipeID)
	if err != nil {
		return nil, err
	}

	// 調味料情報を取得
	seasonings, err := u.recipeRepo.GetRecipeSeasonings(ctx, recipeID)
	if err != nil {
		return nil, err
	}

	// 調理手順を取得
	steps, err := u.recipeRepo.GetRecipeSteps(ctx, recipeID)
	if err != nil {
		return nil, err
	}

	// レシピ詳細DTOを作成
	recipeDetail := &dto.RecipeDetail{
		ID:            recipe.ID,
		RecipeID:      recipe.RecipeID,
		Name:          recipe.Name,
		AuthorComment: recipe.AuthorComment,
		CookTime:      recipe.CookTime,
		Calories:      recipe.Calories,
		TotalPrice:    recipe.TotalPrice,
		CookingPoint:  recipe.CookingPoint,
		ImageURL:      recipe.ImageURL,
		CreatedAt:     recipe.CreatedAt,
		UpdatedAt:     recipe.UpdatedAt,
	}

	// 材料DTOを作成
	var ingredientDTOs []*dto.RecipeIngredient
	for _, ingredient := range ingredients {
		ingredientDTO := &dto.RecipeIngredient{
			ID:           ingredient.ID,
			Name:         ingredient.Name,
			DisplayOrder: ingredient.DisplayOrder,
			AmountText:   ingredient.AmountText,
		}
		ingredientDTOs = append(ingredientDTOs, ingredientDTO)
	}

	// 調味料DTOを作成
	var seasoningDTOs []*dto.RecipeSeasoning
	for _, seasoning := range seasonings {
		seasoningDTO := &dto.RecipeSeasoning{
			ID:           seasoning.ID,
			Name:         seasoning.Name,
			DisplayOrder: seasoning.DisplayOrder,
			AmountText:   seasoning.AmountText,
		}
		seasoningDTOs = append(seasoningDTOs, seasoningDTO)
	}

	// 調理手順DTOを作成
	var stepDTOs []*dto.RecipeStep
	for _, step := range steps {
		stepDTO := &dto.RecipeStep{
			ID:          step.ID,
			Instruction: step.Instruction,
			StepNumber:  step.StepNumber,
		}
		stepDTOs = append(stepDTOs, stepDTO)
	}

	return &dto.GetRecipeDetailResponse{
		Recipe:      recipeDetail,
		Ingredients: ingredientDTOs,
		Seasonings:  seasoningDTOs,
		Steps:       stepDTOs,
	}, nil
}

func (u *RecipeUsecase) GetRecipesByImage(ctx context.Context, req *dto.GetRecipesByImageRequest) (*dto.GetRecipesByImageResponse, error) {
	// 1. 画像から食材を取得
	analysisResult, err := u.openAIService.GetIngredientsFromImage(req.ImageBase64)
	if err != nil {
		return nil, err
	}

	// 2. 解析結果から食材名を抽出（カンマ区切りを分割）
	ingredientNames := strings.Split(analysisResult, ",")
	for i, name := range ingredientNames {
		ingredientNames[i] = strings.TrimSpace(name)
	}

	// 3. 食材名でレシピを検索
	recipes, err := u.recipeRepo.SearchRecipesByIngredients(ctx, ingredientNames)
	if err != nil {
		return nil, err
	}

	// 4. 各レシピの詳細情報を取得
	var searchResults []*dto.RecipeResult
	for _, recipe := range recipes {
		// 食材情報を取得
		ingredients, err := u.recipeRepo.GetRecipeIngredients(ctx, recipe.RecipeID)
		if err != nil {
			return nil, err
		}

		// 調味料情報を取得
		seasonings, err := u.recipeRepo.GetRecipeSeasonings(ctx, recipe.RecipeID)
		if err != nil {
			return nil, err
		}

		// 食材名のリストを作成
		var ingredientNames []string
		for _, ingredient := range ingredients {
			ingredientNames = append(ingredientNames, ingredient.Name)
		}

		// 調味料名のリストを作成
		var seasoningNames []string
		for _, seasoning := range seasonings {
			seasoningNames = append(seasoningNames, seasoning.Name)
		}

		searchResult := &dto.RecipeResult{
			RecipeID:    recipe.RecipeID,
			Name:        recipe.Name,
			CookTime:    recipe.CookTime,
			Calories:    recipe.Calories,
			ImageURL:    recipe.ImageURL,
			Ingredients: ingredientNames,
			Seasonings:  seasoningNames,
		}
		searchResults = append(searchResults, searchResult)
	}

	return &dto.GetRecipesByImageResponse{
		Recipes: searchResults,
	}, nil
}

func (u *RecipeUsecase) SaveRecipe(ctx context.Context, req *dto.SaveRecipeRequest, userID string) (*dto.SaveRecipeResponse, error) {
	// 1. 既に保存されているかチェック
	existingSavedRecipe, err := u.recipeRepo.GetSavedRecipeByUserAndRecipe(ctx, userID, req.RecipeID)
	if err != nil {
		return nil, err
	}
	if existingSavedRecipe != nil {
		return nil, errors.New("recipe is already saved by this user")
	}

	// 2. レシピが存在するかチェック
	recipe, err := u.recipeRepo.GetRecipeByID(ctx, req.RecipeID)
	if err != nil {
		return nil, err
	}
	if recipe == nil {
		return nil, errors.New("recipe not found")
	}

	// 3. 保存レシピエンティティを作成
	savedRecipe := &entity.SavedRecipe{
		SavedRecipeID: uuid.New().String(),
		RecipeID:      req.RecipeID,
		UserID:        userID,
		SavedAt:       time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 4. データベースに保存
	err = u.recipeRepo.SaveRecipe(ctx, savedRecipe)
	if err != nil {
		return nil, err
	}

	return &dto.SaveRecipeResponse{
		SavedRecipeID: savedRecipe.SavedRecipeID,
		RecipeID:      savedRecipe.RecipeID,
		UserID:        savedRecipe.UserID,
		SavedAt:       savedRecipe.SavedAt,
		Message:       "Recipe saved successfully",
	}, nil
}

func (u *RecipeUsecase) GetSavedRecipes(ctx context.Context, userID string) (*dto.GetRecipesByImageResponse, error) {
	// 1. ユーザーが保存したレシピ一覧を取得
	savedRecipes, err := u.recipeRepo.GetSavedRecipesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. 各保存レシピの詳細情報を取得
	var recipeResults []*dto.RecipeResult
	for _, savedRecipe := range savedRecipes {
		// レシピ詳細情報を取得
		recipe, err := u.recipeRepo.GetRecipeByID(ctx, savedRecipe.RecipeID)
		if err != nil {
			return nil, err
		}
		if recipe == nil {
			// レシピが削除されている場合はスキップ
			continue
		}

		// 材料情報を取得
		ingredients, err := u.recipeRepo.GetRecipeIngredients(ctx, savedRecipe.RecipeID)
		if err != nil {
			return nil, err
		}

		// 調味料情報を取得
		seasonings, err := u.recipeRepo.GetRecipeSeasonings(ctx, savedRecipe.RecipeID)
		if err != nil {
			return nil, err
		}

		// 材料名のリストを作成
		var ingredientNames []string
		for _, ingredient := range ingredients {
			ingredientNames = append(ingredientNames, ingredient.Name)
		}

		// 調味料名のリストを作成
		var seasoningNames []string
		for _, seasoning := range seasonings {
			seasoningNames = append(seasoningNames, seasoning.Name)
		}

		// レシピ詳細DTOを作成
		recipeResult := &dto.RecipeResult{
			RecipeID:    recipe.RecipeID,
			Name:        recipe.Name,
			CookTime:    recipe.CookTime,
			Calories:    recipe.Calories,
			ImageURL:    recipe.ImageURL,
			Ingredients: ingredientNames,
			Seasonings:  seasoningNames,
		}

		recipeResults = append(recipeResults, recipeResult)
	}

	return &dto.GetRecipesByImageResponse{
		Recipes: recipeResults,
	}, nil
}
