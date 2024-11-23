/*
Overview
---------
This application is designed to serve the Tidy Notes Android application.
I started this project to help my friends who are still struggling to manage the files provided by our professors.
To address this issue, I created this application. For storage,
I am using Google Drive to manage the files,
and I am hosting the application on Render's free tier.

This application utilizes a JWT authentication strategy without expiration and refresh token mechanisms
for simplicity on Render's free tier.
Despite the lack of expiration,
the system validates tokens to ensure they are not tampered with or invalid.
*/

package main

import (
	"backend/controller"
	"backend/middleware"
	"backend/models"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Initializing the Go Dot Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connecting to the Database.
	models.ConnectToDB()
}

func main() {
	// Setting the gin framework to run on Release Mode.
	gin.SetMode(gin.ReleaseMode)

	// Initializing the GIN Frameword.
	app := gin.Default()

	// Creating the sub-routes.
	auth := app.Group("/auth")
	store := app.Group("/store")

	// JWT Authentication middleware for the store sub-route.
	store.Use(middleware.ValidateJWT)

	// Store
	store.GET("/get", controller.GetFiles)
	store.POST("/add", controller.AddFile)
	store.POST("/delete", controller.DeleteFile)

	// Users
	auth.POST("/login", controller.UserLogin)
	auth.POST("/signup", controller.UserSignUp)
	auth.GET("/validate", middleware.ValidateJWT, controller.ValidateUser)
	auth.GET("/logout", middleware.ValidateJWT, controller.Logout)

	// Getting the port from the .env file.
	PORT := os.Getenv("PORT")
	IP := os.Getenv("IP")
	app.Run(fmt.Sprintf("%s:%s", IP, PORT))
}
