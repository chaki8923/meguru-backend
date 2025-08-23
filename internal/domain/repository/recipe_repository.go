package repository

import (
	"context"
	"meguru-backend/internal/domain/entity"
)

type RecipeRepository interface {
	GetRecipeByID(ctx context.Context, recipeID string) (*entity.Recipe, error)
	GetRecipeIngredients(ctx context.Context, recipeID string) ([]*entity.RecipeIngredient, error)
	GetRecipeSeasonings(ctx context.Context, recipeID string) ([]*entity.RecipeSeasoning, error)
	GetRecipeSteps(ctx context.Context, recipeID string) ([]*entity.RecipeStep, error)
	SearchRecipesByIngredients(ctx context.Context, ingredientNames []string) ([]*entity.Recipe, error)
	SaveRecipe(ctx context.Context, savedRecipe *entity.SavedRecipe) error
	GetSavedRecipeByUserAndRecipe(ctx context.Context, userID, recipeID string) (*entity.SavedRecipe, error)
}
