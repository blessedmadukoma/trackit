package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/blessedmadukoma/trackit-chima/models"
	jwt "github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// Date Formatting
const (
	DDMMYYYYhhmmss = "01/02/2006 15:04:05"
	DDMMYYYY = "01/02/2006"
	hhmm = "15:04"
)

// Dependency injection
type handler struct {
	DB *gorm.DB
}

type Handler = handler

func New(db *gorm.DB) handler {
	return handler{db}
}

func (h handler) checkExistingUser(input *models.User) error {
	emailResult := h.DB.Where("email=?", input.Email).Find(&input)
	phoneResult := h.DB.Where("mobile=?", input.Mobile).Find(&input)
	// usernameResult := h.DB.Where("username=?", input.Username).Find(&input)
	// if emailResult.Error != nil || phoneResult.Error != nil || usernameResult.Error != nil {
	if emailResult.Error != nil || phoneResult.Error != nil {
		err := errors.New("user already exists")
		return err
	}

	return nil
}

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("x-access-token") //Grab the token from the header

		header = strings.TrimSpace(header)

		if header == "" {
			//Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			err := models.ErrorResponse{
				Status:  http.StatusForbidden,
				Message: "Missing auth token",
			}
			json.NewEncoder(w).Encode(err)
			return
		}

		tk := &models.Claims{}

		val := GetJWT()
		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			fmt.Println(val)
			return val, nil
		})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusForbidden)
			err := models.ErrorResponse{
				Status:  http.StatusForbidden,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(err)
			return
		}

		var contextValue = "user"

		ctx := context.WithValue(r.Context(), contextValue, tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func HandleError(input error) models.ErrorResponse {
	err := models.ErrorResponse{
		Message: input.Error(),
		Status:  http.StatusBadRequest,
	}
	return err
}
