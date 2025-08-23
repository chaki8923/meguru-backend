package dto

type SaveRecipeRequest struct {
	RecipeID string `json:"recipe_id" binding:"required"`
}
