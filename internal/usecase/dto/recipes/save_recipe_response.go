package dto

import "time"

type SaveRecipeResponse struct {
	SavedRecipeID string    `json:"saved_recipe_id"`
	RecipeID      string    `json:"recipe_id"`
	UserID        string    `json:"user_id"`
	SavedAt       time.Time `json:"saved_at"`
	Message       string    `json:"message"`
}
