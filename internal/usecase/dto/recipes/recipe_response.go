package dto

import (
	"time"
)

type GetRecipeDetailResponse struct {
	RecipeID    string   `json:"recipe_id"`
	Name        string   `json:"name"`
	CookTime    int      `json:"cook_time"`
	Calories    int      `json:"calories"`
	Ingredients []string `json:"ingredients"`
	Seasonings  []string `json:"seasonings"`
}

type RecipeDetail struct {
	ID            int64     `json:"id"`
	RecipeID      string    `json:"recipe_id"`
	Name          string    `json:"name"`
	AuthorComment string    `json:"author_comment"`
	CookTime      int       `json:"cook_time"`
	Calories      int       `json:"calories"`
	TotalPrice    int       `json:"total_price"`
	CookingPoint  string    `json:"cooking_point"`
	ImageURL      string    `json:"image_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RecipeIngredient struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"display_order"`
	AmountText   string `json:"amount_text"`
}

type RecipeSeasoning struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"display_order"`
	AmountText   string `json:"amount_text"`
}

type RecipeStep struct {
	ID          int64  `json:"id"`
	Instruction string `json:"instruction"`
	StepNumber  int    `json:"step_number"`
}
