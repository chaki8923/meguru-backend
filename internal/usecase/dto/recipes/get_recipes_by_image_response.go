package dto

type GetRecipesByImageResponse struct {
	Recipes []*RecipeResult `json:"recipes"`
}

type RecipeResult struct {
	RecipeID    string   `json:"recipe_id"`
	Name        string   `json:"name"`
	CookTime    int      `json:"cook_time"`
	Calories    int      `json:"calories"`
	ImageURL    string   `json:"image_url,omitempty"`
	Ingredients []string `json:"ingredients"`
	Seasonings  []string `json:"seasonings"`
}
