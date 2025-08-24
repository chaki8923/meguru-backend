package database

import (
	"context"
	"database/sql"
	"fmt"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"

	"github.com/lib/pq"
)

type RecipeRepositoryImpl struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) repository.RecipeRepository {
	return &RecipeRepositoryImpl{db: db}
}

func (r *RecipeRepositoryImpl) GetRecipeByID(ctx context.Context, recipeID string) (*entity.Recipe, error) {
	query := `
		SELECT id, recipe_id, name, author_comment, cook_time, calories, total_price, cooking_point, image_url, created_at, updated_at
		FROM recipes
		WHERE recipe_id = $1 AND deleted_at IS NULL
	`

	recipe := &entity.Recipe{}
	err := r.db.QueryRowContext(ctx, query, recipeID).Scan(
		&recipe.ID, &recipe.RecipeID, &recipe.Name, &recipe.AuthorComment,
		&recipe.CookTime, &recipe.Calories, &recipe.TotalPrice, &recipe.CookingPoint,
		&recipe.ImageURL, &recipe.CreatedAt, &recipe.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe by ID: %w", err)
	}

	return recipe, nil
}

func (r *RecipeRepositoryImpl) GetRecipeIngredients(ctx context.Context, recipeID string) ([]*entity.RecipeIngredient, error) {
	query := `
		SELECT id, recipe_ingredient_id, recipe_id, name, display_order, amount_text, created_at, updated_at
		FROM recipe_ingredients
		WHERE recipe_id = $1 AND deleted_at IS NULL
		ORDER BY display_order
	`

	rows, err := r.db.QueryContext(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []*entity.RecipeIngredient
	for rows.Next() {
		ingredient := &entity.RecipeIngredient{}
		err := rows.Scan(
			&ingredient.ID, &ingredient.RecipeIngredientID, &ingredient.RecipeID,
			&ingredient.Name, &ingredient.DisplayOrder, &ingredient.AmountText,
			&ingredient.CreatedAt, &ingredient.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe ingredient: %w", err)
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func (r *RecipeRepositoryImpl) GetRecipeSeasonings(ctx context.Context, recipeID string) ([]*entity.RecipeSeasoning, error) {
	query := `
		SELECT id, recipe_seasoning_id, recipe_id, name, display_order, amount_text, created_at, updated_at
		FROM recipe_seasonings
		WHERE recipe_id = $1 AND deleted_at IS NULL
		ORDER BY display_order
	`

	rows, err := r.db.QueryContext(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe seasonings: %w", err)
	}
	defer rows.Close()

	var seasonings []*entity.RecipeSeasoning
	for rows.Next() {
		seasoning := &entity.RecipeSeasoning{}
		err := rows.Scan(
			&seasoning.ID, &seasoning.RecipeSeasoningID, &seasoning.RecipeID,
			&seasoning.Name, &seasoning.DisplayOrder, &seasoning.AmountText,
			&seasoning.CreatedAt, &seasoning.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe seasoning: %w", err)
		}
		seasonings = append(seasonings, seasoning)
	}

	return seasonings, nil
}

func (r *RecipeRepositoryImpl) GetRecipeSteps(ctx context.Context, recipeID string) ([]*entity.RecipeStep, error) {
	query := `
		SELECT id, recipe_step_id, recipe_id, instruction, step_number, created_at, updated_at
		FROM recipe_steps
		WHERE recipe_id = $1 AND deleted_at IS NULL
		ORDER BY step_number
	`

	rows, err := r.db.QueryContext(ctx, query, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe steps: %w", err)
	}
	defer rows.Close()

	var steps []*entity.RecipeStep
	for rows.Next() {
		step := &entity.RecipeStep{}
		err := rows.Scan(
			&step.ID, &step.RecipeStepID, &step.RecipeID,
			&step.Instruction, &step.StepNumber,
			&step.CreatedAt, &step.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe step: %w", err)
		}
		steps = append(steps, step)
	}

	return steps, nil
}

func (r *RecipeRepositoryImpl) SearchRecipesByIngredients(ctx context.Context, ingredientNames []string) ([]*entity.Recipe, error) {
	if len(ingredientNames) == 0 {
		return []*entity.Recipe{}, nil
	}

	query := `
		SELECT DISTINCT r.id, r.recipe_id, r.name, r.author_comment, r.cook_time, r.calories, r.total_price, r.cooking_point, r.image_data, r.created_at, r.updated_at
		FROM recipes r
		INNER JOIN recipe_ingredients ri ON r.recipe_id = ri.recipe_id
		WHERE ri.name = ANY($1) AND r.deleted_at IS NULL AND ri.deleted_at IS NULL
		ORDER BY r.name
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(ingredientNames))
	if err != nil {
		return nil, fmt.Errorf("failed to search recipes by ingredients: %w", err)
	}
	defer rows.Close()

	var recipes []*entity.Recipe
	for rows.Next() {
		recipe := &entity.Recipe{}
		err := rows.Scan(
			&recipe.ID, &recipe.RecipeID, &recipe.Name, &recipe.AuthorComment,
			&recipe.CookTime, &recipe.Calories, &recipe.TotalPrice, &recipe.CookingPoint,
			&recipe.ImageURL, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (r *RecipeRepositoryImpl) SaveRecipe(ctx context.Context, savedRecipe *entity.SavedRecipe) error {
	query := `
		INSERT INTO saved_recipes (saved_recipe_id, recipe_id, user_id, saved_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		savedRecipe.SavedRecipeID,
		savedRecipe.RecipeID,
		savedRecipe.UserID,
		savedRecipe.SavedAt,
		savedRecipe.CreatedAt,
		savedRecipe.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save recipe: %w", err)
	}

	return nil
}

func (r *RecipeRepositoryImpl) GetSavedRecipeByUserAndRecipe(ctx context.Context, userID, recipeID string) (*entity.SavedRecipe, error) {
	query := `
		SELECT id, saved_recipe_id, recipe_id, user_id, saved_at, created_at, updated_at
		FROM saved_recipes
		WHERE user_id = $1 AND recipe_id = $2 AND deleted_at IS NULL
	`

	savedRecipe := &entity.SavedRecipe{}
	err := r.db.QueryRowContext(ctx, query, userID, recipeID).Scan(
		&savedRecipe.ID,
		&savedRecipe.SavedRecipeID,
		&savedRecipe.RecipeID,
		&savedRecipe.UserID,
		&savedRecipe.SavedAt,
		&savedRecipe.CreatedAt,
		&savedRecipe.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get saved recipe: %w", err)
	}

	return savedRecipe, nil
}

func (r *RecipeRepositoryImpl) GetSavedRecipesByUser(ctx context.Context, userID string) ([]*entity.SavedRecipe, error) {
	query := `
		SELECT id, saved_recipe_id, recipe_id, user_id, saved_at, created_at, updated_at
		FROM saved_recipes
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY saved_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get saved recipes by user: %w", err)
	}
	defer rows.Close()

	var savedRecipes []*entity.SavedRecipe
	for rows.Next() {
		savedRecipe := &entity.SavedRecipe{}
		err := rows.Scan(
			&savedRecipe.ID,
			&savedRecipe.SavedRecipeID,
			&savedRecipe.RecipeID,
			&savedRecipe.UserID,
			&savedRecipe.SavedAt,
			&savedRecipe.CreatedAt,
			&savedRecipe.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan saved recipe: %w", err)
		}
		savedRecipes = append(savedRecipes, savedRecipe)
	}

	return savedRecipes, nil
}
