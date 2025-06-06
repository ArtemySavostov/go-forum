package handlers

import (
	"log"
	"net/http"
	"strings"

	"JWT/pkg/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		token := parts[1]
		userID, err := authService.ValidateToken(token)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		log.Printf("Username from token: %s", userID)

		c.Set("userID", userID)

		c.Next()
	}
}
