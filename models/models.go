package models

import (
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `gorm:"unique" json:"email"`
	Username  string `gorm:"unique" json:"username"`
	Mobile    string `gorm:"unique" json:"mobile"`
	Password  string `json:"password"`
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

// Organization struct

type TypeEnum string

const (
	Public  TypeEnum = "Public"
	Private TypeEnum = "Private"
	Secret  TypeEnum = "Secret"
)

// type Organization struct {
// 	Name      string          `json:"name"`
// 	CreatorID uint64          `json:"creatorID"`
// 	Type      TypeEnum        `json:"type" sql:"type:ENUM('Public', 'Private', 'Secret')"`
// 	MembersID pq.GenericArray `gorm:"type:Integer[]"`
// 	Logo      string          `json:"logo"`
// }
