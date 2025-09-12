package controller

import (
	"net/http"

	"meguru-backend/internal/usecase"
	dto "meguru-backend/internal/usecase/dto/recipes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// 認証情報があるかチェック
	userID, exists := ctx.Get("user_id")

	var recipeDetail *dto.GetRecipeDetailResponse
	var err error

	if exists {
		// 認証がある場合は保存状態もチェック
		recipeDetail, err = c.recipeUsecase.GetRecipeDetailWithAuth(ctx.Request.Context(), recipeID, userID.(uuid.UUID).String())
	} else {
		// 認証がない場合は通常の取得
		recipeDetail, err = c.recipeUsecase.GetRecipeDetail(ctx.Request.Context(), recipeID)
	}

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

func (c *RecipeController) GetRecipesByImage(ctx *gin.Context) {
	// JWTトークンからユーザーIDを取得
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req dto.GetRecipesByImageRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.recipeUsecase.GetRecipesByImage(ctx.Request.Context(), &req, userID.(uuid.UUID).String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func (c *RecipeController) SaveRecipe(ctx *gin.Context) {
	// JWTトークンからユーザーIDを取得
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req dto.SaveRecipeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.recipeUsecase.SaveRecipe(ctx.Request.Context(), &req, userID.(uuid.UUID).String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func (c *RecipeController) GetSavedRecipes(ctx *gin.Context) {
	// JWTトークンからユーザーIDを取得
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	result, err := c.recipeUsecase.GetSavedRecipes(ctx.Request.Context(), userID.(uuid.UUID).String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func (c *RecipeController) DeleteSavedRecipe(ctx *gin.Context) {
	// JWTトークンからユーザーIDを取得
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// パスパラメータからrecipe_idを取得
	recipeID := ctx.Param("recipe_id")
	if recipeID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "recipe_id is required"})
		return
	}

	// リクエストDTOを作成
	req := &dto.DeleteSavedRecipeRequest{
		RecipeID: recipeID,
	}

	result, err := c.recipeUsecase.DeleteSavedRecipe(ctx.Request.Context(), req, userID.(uuid.UUID).String())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}
