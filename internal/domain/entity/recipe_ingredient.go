package entity

type RecipeIngredient struct {
	ID           uint64 `db:"id"`
	RecipeID     uint64 `db:"recipe_id"`
	IngredientID uint64 `db:"ingredient_id"`
	Quantity     string `db:"quantity"`
}
