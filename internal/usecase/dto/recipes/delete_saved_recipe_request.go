package dto

type DeleteSavedRecipeRequest struct {
	RecipeID string `json:"recipe_id" binding:"required"`
}
