package dto

import (
	"time"
)

type GetRecipeDetailResponse struct {
	Recipe      *RecipeDetail       `json:"recipe"`
	Ingredients []*RecipeIngredient `json:"ingredients"`
	Seasonings  []*RecipeSeasoning  `json:"seasonings"`
	Steps       []*RecipeStep       `json:"steps"`
	SavedFlg    bool                `json:"saved_flg"`
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
