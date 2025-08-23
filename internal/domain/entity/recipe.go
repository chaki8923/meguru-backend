package entity

import (
	"time"
)

type Recipe struct {
	ID            int64
	RecipeID      string
	Name          string
	AuthorComment string
	CookTime      int
	Calories      int
	TotalPrice    int
	CookingPoint  string
	ImageData     []byte
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type RecipeIngredient struct {
	ID                 int64
	RecipeIngredientID string
	RecipeID           string
	Name               string
	DisplayOrder       int
	AmountText         string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
}

type RecipeSeasoning struct {
	ID                int64
	RecipeSeasoningID string
	RecipeID          string
	Name              string
	DisplayOrder      int
	AmountText        string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

type RecipeStep struct {
	ID           int64
	RecipeStepID string
	RecipeID     string
	Instruction  string
	StepNumber   int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type SavedRecipe struct {
	ID            int64
	SavedRecipeID string
	RecipeID      string
	UserID        string
	SavedAt       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
