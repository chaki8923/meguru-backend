package database

import (
	"context"
	"database/sql"
	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type recipeRepository struct {
	db *sqlx.DB
}

func NewRecipeRepository(db *sql.DB) repository.RecipeRepository {
	return &recipeRepository{db: sqlx.NewDb(db, "postgres")}
}

func (r *recipeRepository) FindRecipesByIngredients(ctx context.Context, ingredientNames []string) ([]entity.Recipe, error) {
	query := fmt.Sprintf(`
		SELECT r.*
		FROM recipes r
		JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		JOIN ingredients i ON ri.ingredient_id = i.id
		WHERE i.name IN (?)
		GROUP BY r.id
		HAVING COUNT(DISTINCT i.name) = %d
	`, len(ingredientNames))

	query, args, err := sqlx.In(query, ingredientNames)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)

	var recipes []entity.Recipe
	err = r.db.SelectContext(ctx, &recipes, query, args...)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (r *recipeRepository) SaveRecipe(ctx context.Context, recipe entity.Recipe, ingredients []entity.Ingredient, steps []entity.Step) (entity.Recipe, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return entity.Recipe{}, err
	}
	defer tx.Rollback()

	// Save recipe
	stmt, err := tx.PrepareNamedContext(ctx, `INSERT INTO recipes (title, description, image_url, cooking_time, servings, cost, total_calories, total_score) VALUES (:title, :description, :image_url, :cooking_time, :servings, :cost, :total_calories, :total_score) RETURNING id`)
	if err != nil {
		return entity.Recipe{}, err
	}
	var recipeID uint64
	if err := stmt.GetContext(ctx, &recipeID, recipe); err != nil {
		return entity.Recipe{}, err
	}
	recipe.ID = recipeID

	// Save ingredients
	for _, ingredient := range ingredients {
		var ingredientID uint64
		err := tx.QueryRowxContext(ctx, `INSERT INTO ingredients (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = $1 RETURNING id`, ingredient.Name).Scan(&ingredientID)
		if err != nil {
			return entity.Recipe{}, err
		}

		// Save recipe_ingredients
		_, err = tx.ExecContext(ctx, `INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity) VALUES ($1, $2, $3)`, recipeID, ingredientID, "") // Quantity is not provided in the input, so it's empty.
		if err != nil {
			return entity.Recipe{}, err
		}
	}

	// Save steps
	for _, step := range steps {
		_, err := tx.NamedExecContext(ctx, `INSERT INTO steps (recipe_id, step_number, instruction) VALUES (:recipe_id, :step_number, :instruction)`, map[string]interface{}{
			"recipe_id":   recipeID,
			"step_number": step.StepNumber,
			"instruction": step.Instruction,
		})
		if err != nil {
			return entity.Recipe{}, err
		}
	}

	return recipe, tx.Commit()
}
