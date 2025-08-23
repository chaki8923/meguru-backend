package controller

import (
	"net/http"

	"meguru-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	userUsecase *usecase.UserUsecase
}

func NewUserController(userUsecase *usecase.UserUsecase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var req usecase.CreateUserRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userUsecase.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (uc *UserController) LoginUser(c *gin.Context) {
	var req usecase.LoginUserRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := uc.userUsecase.LoginUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (uc *UserController) ForgotPassword(c *gin.Context) {
	var req usecase.ForgotPasswordRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := uc.userUsecase.ForgotPassword(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (uc *UserController) ResetPassword(c *gin.Context) {
	var req usecase.ResetPasswordRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := uc.userUsecase.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (uc *UserController) GetProfile(c *gin.Context) {
	// JWTトークンからユーザーIDを取得
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なユーザーIDです"})
		return
	}

	response, err := uc.userUsecase.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	// JWTトークンからユーザーIDを取得
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なユーザーIDです"})
		return
	}

	var req usecase.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := uc.userUsecase.UpdateUserProfile(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
} 