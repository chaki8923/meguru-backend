package entity

import "time"

type Step struct {
	ID          uint64    `db:"id"`
	RecipeID    uint64    `db:"recipe_id"`
	StepNumber  int       `db:"step_number"`
	Instruction string    `db:"instruction"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
