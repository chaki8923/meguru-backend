package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"meguru-backend/internal/presentation/responses"
	"meguru-backend/internal/usecase"
	dto "meguru-backend/internal/usecase/dto/users"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (uc *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := uc.userUsecase.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (uc *UserHandler) Signin(c *gin.Context) {
	var req dto.SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := uc.userUsecase.Signin(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			responses.HTTP401(c, gin.H{"error": err.Error()})
			return
		}

		responses.HTTP500(c, gin.H{"error": "internal server error"})
		return
	}

	responses.HTTP201(c, gin.H{
		"token": resp.Token,
		"user":  resp.User,
	})
}

func (uc *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("user_id")

	user, err := uc.userUsecase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
