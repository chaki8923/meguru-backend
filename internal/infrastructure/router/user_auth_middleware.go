package router

import (
	"net/http"
	"strings"

	"meguru-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func UserAuthMiddleware(jwtService *usecase.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
			ctx.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		// ユーザーIDをコンテキストに設定
		ctx.Set("user_id", claims.UserID.String())
		ctx.Set("user_email", claims.Email)
		ctx.Next()
	}
}
