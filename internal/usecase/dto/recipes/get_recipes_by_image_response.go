package dto

import "time"

type GetRecipesByImageResponse struct {
	ExtractedIngredients []string        `json:"extracted_ingredients"`
	LowCalorieRecipes    []*RecipeResult `json:"low_calorie_recipes"`
	LowPriceRecipes      []*RecipeResult `json:"low_price_recipes"`
	QuickCookRecipes     []*RecipeResult `json:"quick_cook_recipes"`
	AIRecommendedRecipes []*RecipeResult `json:"ai_recommended_recipes"`
}

type RecipeResult struct {
	RecipeID    string    `json:"recipe_id"`
	Name        string    `json:"name"`
	CookTime    int       `json:"cook_time"`
	Calories    int       `json:"calories"`
	TotalPrice  int       `json:"total_price"`
	ImageURL    string    `json:"image_url,omitempty"`
	Ingredients []string  `json:"ingredients"`
	Seasonings  []string  `json:"seasonings"`
	SavedFlg    bool      `json:"saved_flg"`
	CreatedAt   time.Time `json:"created_at"`
}
