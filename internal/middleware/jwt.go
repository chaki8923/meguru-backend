package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecretKey = []byte("your-very-secret-key")

func GenerateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecretKey)
}

func ValidateJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !checkAuthHeaderPresence(c, authHeader) {
			return
		}

		tokenString := extractBearerToken(c, authHeader)
		if tokenString == "" {
			return
		}

		token := parseAndValidateToken(c, tokenString)
		if token == nil {
			return
		}

		if !checkTokenClaims(c, token) {
			return
		}

		c.Next()
	}
}

func checkAuthHeaderPresence(c *gin.Context, authHeader string) bool {
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		c.Abort()
		return false
	}
	return true
}

func checkTokenClaims(c *gin.Context, token *jwt.Token) bool {
	_, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		c.Abort()
		return false
	}
	return true
}

func extractBearerToken(c *gin.Context, authHeader string) string {
	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer {token}"})
		c.Abort()
		return ""
	}
	return parts[1]
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)

	if !ok {
		return nil, jwt.ErrTokenMalformed
	}
	return jwtSecretKey, nil
}

func parseAndValidateToken(c *gin.Context, tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, keyFunc)

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		c.Abort()
		return nil
	}

	return token
}
