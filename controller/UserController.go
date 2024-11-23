package controller

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func UserSignUp(c *gin.Context) {
	// Declaring the required variables
	var user models.Users
	var hashedPassword string

	// Decode the request body
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON format",
		})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.UserPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Password hashing failed",
		})
		return
	}
	user.UserPassword = hashedPassword

	// Insert the record into the database
	if err := models.DB.Create(&user).Error; err != nil {
		// Handle unique constraint violation (PostgreSQL error code 23505)
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Username or Email already exists",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
		}
		return
	}

	// Successful signup
	c.JSON(http.StatusOK, gin.H{
		"message": "Signup Successful",
	})
}

func UserLogin(c *gin.Context) {
	var user models.Users
	var expectedUser models.Users

	// Decode the request body
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON format",
		})
		return
	}

	// Fetch the user from the database
	if err := models.DB.Where("useremail = ?", user.UserEmail).First(&expectedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database Error",
			})
		}
		return
	}

	// Validate password
	if utils.ComparePasswordAndHashedPassword(expectedUser.UserPassword, user.UserPassword) {
		if !expectedUser.IsLoggedIn {
			tokenString, err := utils.GenerateToken(expectedUser.UserID, expectedUser.UserName)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to generate token",
				})
				return
			}

			expectedUser.IsLoggedIn = true
			models.DB.Save(&expectedUser)

			c.JSON(http.StatusOK, gin.H{
				"message": "Login successful",
				"token":   tokenString,
			})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Account logged in on another device",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "Incorrect Password",
	})
}

func Logout(c *gin.Context) {
	var expectedUser models.Users

	userID := c.MustGet("userid")
	username := c.MustGet("username")

	// Fetch the user from the database
	if err := models.DB.Where("userid = ? AND username = ?", userID, username).First(&expectedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database Error",
			})
		}
		return
	}

	// Update `isLoggedIn` status
	if expectedUser.IsLoggedIn {
		expectedUser.IsLoggedIn = false
		models.DB.Save(&expectedUser)
		c.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "Logout unsuccessful",
	})
}

func ValidateUser(c *gin.Context) {
	var expectedUser models.Users

	userID := c.MustGet("userid")
	username := c.MustGet("username")

	// Fetch the user from the database
	if err := models.DB.Where("userid = ? AND username = ?", userID, username).First(&expectedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database Error",
			})
		}
		return
	}

	if expectedUser.IsLoggedIn {
		c.JSON(http.StatusOK, gin.H{
			"message": "Validation successful",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "Validation unsuccessful",
	})
}
