package controller

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func UserSignUp(c *gin.Context) {
	// Declaring the requird variables
	var user models.Users
	var hashedPassword string

	// After the receiving the POST request, we are trying to decode the body of the request.
	// Decoding the body of the request to the type of `Users`struct and storing it in (user).
	if json.NewDecoder(c.Request.Body).Decode(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid json format",
		})
		return
	}

	// trying to hash the password using the bycrypt package.
	hashedPassword, err := utils.HashPassword(user.UserPassword)

	// if any error occured during hashing the password we responded the client with InternalServerError.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Password hashing failed",
		})
		return
	}

	// If hashing is successful, change the (UserPassword) field in the (user).
	user.UserPassword = hashedPassword

	user.UserID = 0

	// Inserting the record to the Database.
	if err := models.DB.Create(&user).Error; err != nil {
		// Checks for violation of unique key constraint.
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
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

	// if there is no error. Then, we send the client with signup successful message.
	// with a status code of (200 OK).
	c.JSON(http.StatusOK, gin.H{
		"message": "Signup Successful",
	})
}

func UserLogin(c *gin.Context) {
	// Declaring the required variables.
	var user models.Users
	var expected_user models.Users

	// After the receiving the POST request, we are trying to decode the body of the request.
	// Decoding the body of the request to the type of `Users`struct and storing it in (user).
	if json.NewDecoder(c.Request.Body).Decode(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid json format",
		})
		return
	}

	// Fetching the user record from the database.
	err := models.DB.Where("useremail = ?", user.UserEmail).First(&expected_user).Error

	if err != nil {
		// Check for no records found error.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found",
			})
		} else {
			// For other errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database Error",
			})
		}
		return
	}

	// Comparing the string password (user.UserPassword) with the hashed password (expected_user.UserPassword).
	if utils.ComparePasswordAndHashedPassword(expected_user.UserPassword, user.UserPassword) {
		// Check if the user account is already logged in another device.
		if !expected_user.IsLoggedIn {
			// Generating the JWT Token.
			tokenString, err := utils.GenerateToken(expected_user.UserID, expected_user.UserName)
			// Checking for any error while generating the token.
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to generate token",
				})
			}
			// Update the `isLoggedIn` column as true.
			expected_user.IsLoggedIn = true
			models.DB.Save(&expected_user)

			// After successfully generating the JWT token.
			// we send that token to the client with login successful message.
			c.JSON(http.StatusOK, gin.H{
				"message": "login successful",
				"token":   tokenString,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Account logged in another device",
			})
			return
		}

	} else {
		// When hash values of the passwords are different.
		// Then we notify the client as login failed message.
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Incorrect Password",
		})
		return
	}
}

func Logout(c *gin.Context) {
	// Declaring the necessary variable.
	var expected_user models.Users

	// Getting the contents from the middleware.
	userid := c.MustGet("userid")
	username := c.MustGet("username")

	// Fetching the user record from the database.
	err := models.DB.Where("userid = ? AND username = ?", userid, username).First(&expected_user).Error

	if err != nil {
		// Check for no records found error.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found",
			})
		} else {
			// For other errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database Error",
			})
		}
		return
	}

	// updating the column
	if expected_user.IsLoggedIn {
		expected_user.IsLoggedIn = false
		models.DB.Save(&expected_user)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "logout Unsuccessful",
		})
	}

	// Logout success response.
	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}

func ValidateUser(c *gin.Context) {
	var expected_user models.Users

	// Getting the contents from the middleware.
	userid := c.MustGet("userid")
	username := c.MustGet("username")

	// Fetching the user record from the database.
	err := models.DB.Where("userid = ? AND username = ?", userid, username).First(&expected_user).Error

	if err != nil {
		// Check for no records found error.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found",
			})
		} else {
			// For other errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database Error",
			})
		}
		return
	}

	if expected_user.IsLoggedIn {
		c.JSON(http.StatusOK, gin.H{
			"message": "validation successful",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "validation Unsuccessful",
	})
}
