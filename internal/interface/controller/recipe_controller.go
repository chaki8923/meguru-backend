package controller

import (
	"net/http"

	"meguru-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type RecipeController struct {
	recipeUsecase *usecase.RecipeUsecase
}

func NewRecipeController(recipeUsecase *usecase.RecipeUsecase) *RecipeController {
	return &RecipeController{
		recipeUsecase: recipeUsecase,
	}
}

func (c *RecipeController) GetRecipeDetail(ctx *gin.Context) {
	recipeID := ctx.Param("recipe_id")

	recipeDetail, err := c.recipeUsecase.GetRecipeDetail(ctx.Request.Context(), recipeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if recipeDetail == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": recipeDetail})
}
