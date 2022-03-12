package models

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `gorm:"unique" json:"email"`
	// Username  string `gorm:"unique" json:"username"`
	Mobile   string `gorm:"unique" json:"mobile"`
	Password string `json:"password"`
}

type Expense struct {
	gorm.Model
	Amount         float64   `json:"amount"`
	Description    string    `json:"description"`
	Date_purchased time.Time `json:"date_purchased"`
	Category       string    `json:"category"`
	UserID         uint
	User           User
	// User        User  `gorm:"foreignkey:UserID;association_foreignkey:ID" json:"UserID"`
}

type Budget struct {
	gorm.Model
	Budget_name string    `json:"budget_name"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate"`
	EndDate     string    `json:"endDate"`
	UserID      uint
	User        User
	// User        User  `gorm:"foreignkey:UserID;association_foreignkey:ID" json:"UserID"`
}

type Transactions struct {
	gorm.Model
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	UserID   uint
	User     User
	// User        User  `gorm:"foreignkey:UserID;association_foreignkey:ID" json:"UserID"`
}

type Income struct {
	gorm.Model
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
	UserID uint
	User   User
}

type Account struct {
	gorm.Model
	Amount float64 `json:"amount"`
	UserID uint
	User   User
}

type ErrorResponse struct {
	Status  int
	Message string
}

type Claims struct {
	User User
	jwt.StandardClaims
}

type LoginUser struct {
	Email    string
	Password string
}

type ResetUser struct {
	Email        string
	New_password string
}

// type TypeEnum string

// const (
// 	Public  TypeEnum = "Public"
// 	Private TypeEnum = "Private"
// 	Secret  TypeEnum = "Secret"
// )
