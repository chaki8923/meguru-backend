package usecase

import (
	"context"
	"errors"
	"fmt"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"meguru-backend/internal/infrastructure/service"
	dto "meguru-backend/internal/usecase/dto/recipes"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RecipeUsecase struct {
	recipeRepo    repository.RecipeRepository
	userRepo      repository.UserRepository
	openAIService *service.OpenAIService
}

func NewRecipeUsecase(recipeRepo repository.RecipeRepository, userRepo repository.UserRepository, openAIService *service.OpenAIService) *RecipeUsecase {
	return &RecipeUsecase{
		recipeRepo:    recipeRepo,
		userRepo:      userRepo,
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
		SavedFlg:    false, // 認証なしの場合はデフォルトでfalse
	}, nil
}

func (u *RecipeUsecase) GetRecipeDetailWithAuth(ctx context.Context, recipeID, userID string) (*dto.GetRecipeDetailResponse, error) {
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

	// レシピが保存済みかどうかをチェック
	savedRecipe, err := u.recipeRepo.GetSavedRecipeByUserAndRecipe(ctx, userID, recipeID)
	if err != nil {
		return nil, err
	}
	savedFlg := savedRecipe != nil

	return &dto.GetRecipeDetailResponse{
		Recipe:      recipeDetail,
		Ingredients: ingredientDTOs,
		Seasonings:  seasoningDTOs,
		Steps:       stepDTOs,
		SavedFlg:    savedFlg,
	}, nil
}

// shouldExcludeRecipe は、画像から抽出した食材とレシピの食材を比較して、
// 特定の組み合わせを除外するかどうかを判定する
func (u *RecipeUsecase) shouldExcludeRecipe(extractedIngredients []string, recipeIngredients []string) bool {
	// 画像から抽出した食材をマップに変換（高速検索のため）
	extractedMap := make(map[string]bool)
	for _, ingredient := range extractedIngredients {
		extractedMap[ingredient] = true
	}

	// レシピの食材をチェック
	for _, recipeIngredient := range recipeIngredients {
		// トマトの除外ルール
		if extractedMap["トマト"] {
			if strings.Contains(recipeIngredient, "ミニトマト") ||
				strings.Contains(recipeIngredient, "カットトマト") {
				return true
			}
		}

		// 玉ねぎの除外ルール
		if extractedMap["玉ねぎ"] {
			if strings.Contains(recipeIngredient, "紫玉ねぎ（薄切り）") {
				return true
			}
		}
	}

	return false
}

func (u *RecipeUsecase) GetRecipesByImage(ctx context.Context, req *dto.GetRecipesByImageRequest, userID string) (*dto.GetRecipesByImageResponse, error) {
	// 1. 画像から食材を抽出
	ingredientsText, err := u.openAIService.GetIngredientsFromImage(req.ImageBase64)
	if err != nil {
		return nil, err
	}

	// 2. カンマ区切りで分割して食材名のリストを作成
	ingredientNames := strings.Split(ingredientsText, ",")
	for i, name := range ingredientNames {
		ingredientNames[i] = strings.TrimSpace(name)
	}

	// 3. 食材名でレシピを検索
	recipes, err := u.recipeRepo.SearchRecipesByIngredients(ctx, ingredientNames)
	if err != nil {
		return nil, err
	}

	// 4. 各レシピの詳細情報を取得
	var allRecipes []*dto.RecipeResult
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
		var recipeIngredientNames []string
		for _, ingredient := range ingredients {
			recipeIngredientNames = append(recipeIngredientNames, ingredient.Name)
		}

		// 除外ルールをチェック
		if u.shouldExcludeRecipe(ingredientNames, recipeIngredientNames) {
			continue // このレシピを除外
		}

		// 調味料名のリストを作成
		var seasoningNames []string
		for _, seasoning := range seasonings {
			seasoningNames = append(seasoningNames, seasoning.Name)
		}

		// レシピが保存済みかどうかをチェック
		savedRecipe, err := u.recipeRepo.GetSavedRecipeByUserAndRecipe(ctx, userID, recipe.RecipeID)
		if err != nil {
			return nil, err
		}
		savedFlg := savedRecipe != nil

		recipeResult := &dto.RecipeResult{
			RecipeID:    recipe.RecipeID,
			Name:        recipe.Name,
			CookTime:    recipe.CookTime,
			Calories:    recipe.Calories,
			TotalPrice:  recipe.TotalPrice,
			ImageURL:    recipe.ImageURL,
			Ingredients: recipeIngredientNames,
			Seasonings:  seasoningNames,
			SavedFlg:    savedFlg,
		}
		allRecipes = append(allRecipes, recipeResult)
	}

	// 5. カロリー順（昇順）でソート
	lowCalorieRecipes := make([]*dto.RecipeResult, len(allRecipes))
	copy(lowCalorieRecipes, allRecipes)
	sort.Slice(lowCalorieRecipes, func(i, j int) bool {
		return lowCalorieRecipes[i].Calories < lowCalorieRecipes[j].Calories
	})
	if len(lowCalorieRecipes) > 10 {
		lowCalorieRecipes = lowCalorieRecipes[:10]
	}

	// 6. 金額順（昇順）でソート
	lowPriceRecipes := make([]*dto.RecipeResult, len(allRecipes))
	copy(lowPriceRecipes, allRecipes)
	sort.Slice(lowPriceRecipes, func(i, j int) bool {
		return lowPriceRecipes[i].TotalPrice < lowPriceRecipes[j].TotalPrice
	})
	if len(lowPriceRecipes) > 10 {
		lowPriceRecipes = lowPriceRecipes[:10]
	}

	// 7. 調理時間順（昇順）でソート
	quickCookRecipes := make([]*dto.RecipeResult, len(allRecipes))
	copy(quickCookRecipes, allRecipes)
	sort.Slice(quickCookRecipes, func(i, j int) bool {
		return quickCookRecipes[i].CookTime < quickCookRecipes[j].CookTime
	})
	if len(quickCookRecipes) > 10 {
		quickCookRecipes = quickCookRecipes[:10]
	}

	// 8. AIおすすめレシピを取得（完全に独立したランダム選択）
	// 元の全レシピから直接選ぶ（他の3つのカテゴリとは独立）
	var aiRecommendedRecipes []*dto.RecipeResult
	if len(allRecipes) > 0 {
		// 全レシピ名を収集
		var allRecipeNames []string
		for _, recipe := range allRecipes {
			allRecipeNames = append(allRecipeNames, recipe.Name)
		}

		// OpenAI APIでAIおすすめレシピを取得
		recommendedNames, err := u.openAIService.GetRecommendedRecipes(allRecipeNames)
		if err != nil {
			// OpenAI APIエラーの場合は、全レシピからランダムに10件を返す
			if len(allRecipes) > 10 {
				aiRecommendedRecipes = allRecipes[:10]
			} else {
				aiRecommendedRecipes = allRecipes
			}
		} else {
			// AIおすすめレシピ名に基づいてレシピを選択
			recipeMap := make(map[string]*dto.RecipeResult)
			for _, recipe := range allRecipes {
				recipeMap[recipe.Name] = recipe
			}

			for _, recommendedName := range recommendedNames {
				if recipe, exists := recipeMap[recommendedName]; exists {
					aiRecommendedRecipes = append(aiRecommendedRecipes, recipe)
					if len(aiRecommendedRecipes) >= 10 {
						break
					}
				}
			}
		}
	}

	return &dto.GetRecipesByImageResponse{
		ExtractedIngredients: ingredientNames,
		LowCalorieRecipes:    lowCalorieRecipes,
		LowPriceRecipes:      lowPriceRecipes,
		QuickCookRecipes:     quickCookRecipes,
		AIRecommendedRecipes: aiRecommendedRecipes,
	}, nil
}

func (u *RecipeUsecase) SaveRecipe(ctx context.Context, req *dto.SaveRecipeRequest, userID string) (*dto.SaveRecipeResponse, error) {
	// デバッグログ追加
	fmt.Printf("SaveRecipe called with userID: %s, recipeID: %s\n", userID, req.RecipeID)

	// 1. ユーザーが存在するかチェック
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		fmt.Printf("Invalid user ID format: %s\n", userID)
		return nil, errors.New("invalid user ID format")
	}

	user, err := u.userRepo.GetByID(ctx, userUUID)
	if err != nil {
		fmt.Printf("Error getting user by ID: %v\n", err)
		return nil, err
	}
	if user == nil {
		fmt.Printf("User not found: %s\n", userID)
		return nil, errors.New("user not found")
	}

	fmt.Printf("User found: %s\n", user.Email)

	// 2. 既に保存されているかチェック
	existingSavedRecipe, err := u.recipeRepo.GetSavedRecipeByUserAndRecipe(ctx, userID, req.RecipeID)
	if err != nil {
		return nil, err
	}
	if existingSavedRecipe != nil {
		return nil, errors.New("recipe is already saved by this user")
	}

	// 3. レシピが存在するかチェック
	recipe, err := u.recipeRepo.GetRecipeByID(ctx, req.RecipeID)
	if err != nil {
		return nil, err
	}
	if recipe == nil {
		return nil, errors.New("recipe not found")
	}

	// 4. 保存レシピエンティティを作成
	savedRecipe := &entity.SavedRecipe{
		SavedRecipeID: uuid.New().String(),
		RecipeID:      req.RecipeID,
		UserID:        userID,
		SavedAt:       time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 5. データベースに保存
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

func (u *RecipeUsecase) GetSavedRecipes(ctx context.Context, userID string) (*dto.GetSavedRecipesResponse, error) {
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
			TotalPrice:  recipe.TotalPrice,
			ImageURL:    recipe.ImageURL,
			Ingredients: ingredientNames,
			Seasonings:  seasoningNames,
			SavedFlg:    true,                  // 保存済みレシピなのでtrue
			CreatedAt:   savedRecipe.CreatedAt, // 保存日時を使用
		}

		recipeResults = append(recipeResults, recipeResult)
	}

	return &dto.GetSavedRecipesResponse{
		Recipes: recipeResults,
	}, nil
}

func (u *RecipeUsecase) DeleteSavedRecipe(ctx context.Context, req *dto.DeleteSavedRecipeRequest, userID string) (*dto.DeleteSavedRecipeResponse, error) {
	// 1. 保存レシピが存在するかチェック
	savedRecipe, err := u.recipeRepo.GetSavedRecipeByUserAndRecipe(ctx, userID, req.RecipeID)
	if err != nil {
		return nil, err
	}
	if savedRecipe == nil {
		return nil, errors.New("saved recipe not found")
	}

	// 2. 保存レシピを削除（論理削除）
	err = u.recipeRepo.DeleteSavedRecipe(ctx, userID, req.RecipeID)
	if err != nil {
		return nil, err
	}

	return &dto.DeleteSavedRecipeResponse{
		Message: "Recipe removed from saved recipes successfully",
	}, nil
}
