package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

	// Initializing the Database.
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require", dbuser, password, databaseip, databaseport, database)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tidynotes.", // Prefix if needed
			SingularTable: true,         // Use singular table names
		},
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Database connection established successfully.", DB)
}
