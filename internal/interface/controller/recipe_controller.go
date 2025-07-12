package controller

import (
	"net/http"
	"strings"

	"meguru-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type RecipeController struct {
	recipeUsecase *usecase.RecipeUsecase
}

func NewRecipeController(recipeUsecase *usecase.RecipeUsecase) *RecipeController {
	return &RecipeController{recipeUsecase: recipeUsecase}
}

func (c *RecipeController) SuggestRecipes(ctx *gin.Context) {
	ingredientsParam := ctx.Query("ingredients")
	if ingredientsParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ingredients are required"})
		return
	}
	ingredients := strings.Split(ingredientsParam, ",")

	recipes, err := c.recipeUsecase.SuggestRecipes(ctx.Request.Context(), ingredients)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, recipes)
}
