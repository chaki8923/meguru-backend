package router

import (
	"fmt"
	"net/http"
	"strings"

	"meguru-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getTokenFromHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("token is required")
	}

	return token, nil
}

func getStoreIDFromToken(token string) (uuid.UUID, error) {
	var uuidStr string
	if strings.HasPrefix(token, "auth_token_") {
		uuidStr = strings.TrimPrefix(token, "auth_token_")
	} else if strings.HasPrefix(token, "temp_token_") {
		uuidStr = strings.TrimPrefix(token, "temp_token_")
	} else {
		return uuid.Nil, fmt.Errorf("invalid token format")
	}

	storeID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	return storeID, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getTokenFromHeader(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		storeID, err := getStoreIDFromToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("storeID", storeID)
		ctx.Next()
	}
}

// UserAuthMiddleware ユーザー認証用のミドルウェア
func UserAuthMiddleware(jwtService *usecase.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getTokenFromHeader(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID.String())
		ctx.Set("user_email", claims.Email)
		ctx.Next()
	}
}
