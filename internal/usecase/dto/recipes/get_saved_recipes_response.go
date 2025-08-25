package dto

type GetSavedRecipesResponse struct {
	Recipes []*RecipeResult `json:"recipes"`
}
