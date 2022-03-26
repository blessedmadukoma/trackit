package models

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload" // helps to autoload the env

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading env file:", err)
	// }

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

	err = db.AutoMigrate(&User{}, &Expense{}, &Budget{}, &Transactions{}, &Income{}, &Account{})
	if err != nil {
		log.Fatal("Error migrating db:", err)
	}

	return db
}
