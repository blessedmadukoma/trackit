package models

import (
	"fmt"
	"log"
	"os"

	// "github.com/blessedmadukoma/trackit-chima/models"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file:", err)
	}

	// dbESQL := os.Getenv("ELEPHANTSQL_URI")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbTZone := os.Getenv("DB_TIMEZONE")

	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s", dbHost, dbUser, dbPass, dbName, dbPort, dbTZone)
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("DB Connected!!")

	// AutoMigrate the User and Organization structs: comment out if done
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Error migrating db:", err)
	}

	return db
}
