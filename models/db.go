package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	// Getting all the data from Environmental Variables.
	dbuser := os.Getenv("DBUSER")
	password := os.Getenv("PASSWORD")
	database := os.Getenv("DATABASE")
	databaseip := os.Getenv("DBIP")
	databaseport := os.Getenv("DBPORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbuser, password, databaseip, databaseport, database)
	// Initializing the Database.
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
		return
	}

	log.Println("Database connection established successfully.")
}
