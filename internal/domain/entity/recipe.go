package entity

import "time"

type Recipe struct {
	ID          uint64    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	ImageURL    string    `db:"image_url"`
	CookingTime int       `db:"cooking_time"`
	Servings    int       `db:"servings"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
