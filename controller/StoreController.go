package controller

import (
	"backend/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func GetFiles(c *gin.Context) {
	var files []models.Stores

	// Getting the list of files from the Database.
	if err := models.DB.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database Error",
		})
		return
	}

	// Responding the client with the list of files.
	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})

}

func AddFile(c *gin.Context) {
	var file models.Stores

	// Decode JSON from the request body
	if err := json.NewDecoder(c.Request.Body).Decode(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON format",
		})
		return
	}

	// Inserting the record into the Database
	if err := models.DB.Create(&file).Error; err != nil {
		// Check for violation of unique key constraint (MySQL error code 1062)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "File already exists",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
		}
		return
	}

	// Return success response along with the created file object
	c.JSON(http.StatusOK, gin.H{
		"message": "File created successfully",
		"file":    file, // Returning the file object to confirm creation
	})
}

func DeleteFile(c *gin.Context) {
	// Define the File struct with an exported field
	type File struct {
		Filename string `json:"filename"` // Use uppercase "Filename"
	}

	var file File

	// Decode the incoming JSON into the file struct
	if err := json.NewDecoder(c.Request.Body).Decode(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid JSON format",
		})
		return
	}

	// Perform the deletion operation
	result := models.DB.Where("filename = ?", file.Filename).Delete(&models.Stores{})

	// Check if any rows were affected
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to delete the file",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found",
		})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}
