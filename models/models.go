package models

import (
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"gorm.io/gorm"

	_ "github.com/go-playground/validator/v10"
)

type User struct {
	gorm.Model
	Firstname string `json:"firstname" gorm:"not null"`
	Lastname  string `json:"lastname" gorm:"not null"`
	Email     string `gorm:"unique;not null" json:"email"`
	// Username  string `gorm:"unique" json:"username"`
	Mobile   string `gorm:"unique;not null" json:"mobile"`
	Password string `json:"password" gorm:"not null"`
}

type Expense struct {
	gorm.Model
	Amount      float64 `gorm:"not null" json:"amount"`
	Description string  `gorm:"not null" json:"description"`
	// Date_purchased time.Time `json:"date_purchased"`
	Date_purchased string `json:"date_purchased"`
	Category       string `gorm:"not null" json:"category"`
	UserID         uint
	User           User
	// User        User  `gorm:"foreignkey:UserID;association_foreignkey:ID" json:"UserID"`
}

type Budget struct {
	gorm.Model
	Budget_name string  `json:"budget_name" gorm:"not null"`
	Amount      float64 `json:"amount" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	StartDate   string  `json:"startDate" gorm:"not null"`
	EndDate     string  `json:"endDate" gorm:"not null"`
	UserID      uint
	User        User
	// User        User  `gorm:"foreignkey:UserID;association_foreignkey:ID" json:"UserID"`
}

type Transactions struct {
	gorm.Model
	Category string  `json:"category" gorm:"not null"`
	Amount   float64 `json:"amount" gorm:"not null"`
	Date string `json:"date" gorm:"not null"`
	Time string `json:"time" gorm:"not null"`
	UserID   uint
	User     User
	// User        User  `gorm:"foreignkey:UserID;association_foreignkey:ID" json:"UserID"`
}

type Income struct {
	gorm.Model
	Amount float64 `json:"amount" gorm:"not null"`
	// Date   time.Time `json:"date" gorm:"not null"`
	Date   string `json:"date" gorm:"not null"`
	UserID uint
	User   User
}

type Account struct {
	gorm.Model
	Amount float64 `json:"amount" gorm:"not null"`
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
