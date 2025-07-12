package repository

import (
	"context"
	"meguru-backend/internal/domain/entity"
)

type RecipeRepository interface {
	FindRecipesByIngredients(ctx context.Context, ingredientNames []string) ([]entity.Recipe, error)
	SaveRecipe(ctx context.Context, recipe entity.Recipe, ingredients []entity.Ingredient, steps []entity.Step) (entity.Recipe, error)
}
