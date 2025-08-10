package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/abhoy21/expense-tracker/initializers"
	"github.com/abhoy21/expense-tracker/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateAuth(c *gin.Context) {
	var tokenString string

	// Try Authorization header (mobile clients)
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Fallback to cookie (web clients)
	if tokenString == "" {
		cookieToken, err := c.Cookie("Authorization")
		if err == nil {
			tokenString = cookieToken
		}
	}

	// No token found
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authentication token provided"})
		return
	}

	// Parse and validate JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Expiration check
	if exp, ok := claims["exp"].(float64); ok {
		if float64(time.Now().Unix()) > exp {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session expired, please log in again"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token expiration"})
		return
	}

	// Load user from DB
	userID, ok := claims["id"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
		return
	}

	var user models.User
	if err := initializers.DBClient.Select("ID", "Email").First(&user, uint(userID)).Error; err != nil || user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found or unauthorized"})
		return
	}

	// Pass user to next handlers
	c.Set("user", user)
	c.Next()
}
