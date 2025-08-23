package usecase

import (
	"context"
	"meguru-backend/internal/domain/repository"
	dto "meguru-backend/internal/usecase/dto/recipes"
)

type RecipeUsecase struct {
	recipeRepo repository.RecipeRepository
}

func NewRecipeUsecase(recipeRepo repository.RecipeRepository) *RecipeUsecase {
	return &RecipeUsecase{
		recipeRepo: recipeRepo,
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

	// 手順情報を取得
	steps, err := u.recipeRepo.GetRecipeSteps(ctx, recipeID)
	if err != nil {
		return nil, err
	}

	// DTOに変換
	recipeDetail := &dto.RecipeDetail{
		ID:            recipe.ID,
		RecipeID:      recipe.RecipeID,
		Name:          recipe.Name,
		AuthorComment: recipe.AuthorComment,
		CookTime:      recipe.CookTime,
		Calories:      recipe.Calories,
		TotalPrice:    recipe.TotalPrice,
		CookingPoint:  recipe.CookingPoint,
		ImageData:     recipe.ImageData,
		CreatedAt:     recipe.CreatedAt,
		UpdatedAt:     recipe.UpdatedAt,
	}

	// 材料DTOに変換
	var ingredientDTOs []*dto.RecipeIngredient
	for _, ingredient := range ingredients {
		ingredientDTOs = append(ingredientDTOs, &dto.RecipeIngredient{
			ID:           ingredient.ID,
			Name:         ingredient.Name,
			DisplayOrder: ingredient.DisplayOrder,
			AmountText:   ingredient.AmountText,
		})
	}

	// 調味料DTOに変換
	var seasoningDTOs []*dto.RecipeSeasoning
	for _, seasoning := range seasonings {
		seasoningDTOs = append(seasoningDTOs, &dto.RecipeSeasoning{
			ID:           seasoning.ID,
			Name:         seasoning.Name,
			DisplayOrder: seasoning.DisplayOrder,
			AmountText:   seasoning.AmountText,
		})
	}

	// 手順DTOに変換
	var stepDTOs []*dto.RecipeStep
	for _, step := range steps {
		stepDTOs = append(stepDTOs, &dto.RecipeStep{
			ID:          step.ID,
			Instruction: step.Instruction,
			StepNumber:  step.StepNumber,
		})
	}

	return &dto.GetRecipeDetailResponse{
		Recipe:      recipeDetail,
		Ingredients: ingredientDTOs,
		Seasonings:  seasoningDTOs,
		Steps:       stepDTOs,
	}, nil
}
