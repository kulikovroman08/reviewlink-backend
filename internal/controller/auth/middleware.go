package auth

import (
	"fmt"
	jwt2 "github.com/kulikovroman08/reviewlink-backend/pkg/jwt"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims, err := parseToken(tokenStr)
		if err != nil {
			slog.Error("token parsing failed", "error", err)
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("claims", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func parseToken(tokenStr string) (*jwt2.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt2.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	claims, ok := token.Claims.(*jwt2.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}
