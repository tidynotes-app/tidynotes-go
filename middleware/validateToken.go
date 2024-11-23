package middleware

import (
	"backend/models"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// ValidateJWT middleware to check the validity of the JWT token
func ValidateJWT(c *gin.Context) {
	var user models.Users
	jwtkey := os.Getenv("JWTKEY")

	// Extract token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is required"})
		c.Abort()
		return
	}

	// Check if the Authorization header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
		c.Abort()
		return
	}

	// Get the token part from the Authorization header
	tokenString := authHeader[7:] // Remove "Bearer " prefix

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HS256 (HMAC-SHA256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation
		return []byte(jwtkey), nil
	})

	// Check for parsing errors
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		c.Abort()
		return
	}

	// Ensure the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		userID := claims["userid"]
		userName := claims["username"]

		c.Set("userid", userID)
		c.Set("username", userName)

		if models.DB.Where("userid = ? AND username = ?", userID, userName).First(&user).Error != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// User not found
				c.JSON(http.StatusNotFound, gin.H{
					"message": "User not found",
				})
			} else {
				// Other database error
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Database Error",
				})
			}
		}

		// Continue to the next handler
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
		c.Abort()
	}
}
